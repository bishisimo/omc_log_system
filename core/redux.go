/*
@author '彼时思默'
@time 2020/4/29 上午9:55
@describe: 发布订阅者统一管理
*/
package core

import (
	"encoding/json"
	"github.com/roylee0704/gron"
	"github.com/roylee0704/gron/xtime"
	"github.com/sirupsen/logrus"
	"os"
)

var Redux *redux

func init() {
	var isDrop = true
	if os.Getenv("IsRetainLocal") == "true" {
		isDrop = false
	}
	Redux = NewRedux()
	Redux.TimerDistribution(isDrop)
}

type redux struct {
	Local     *LocalSub //map[string]*LocalSub 本地储存订阅者通道
	SubsFlow  SubMap    //map[string]*Sub 流数据订阅者通道
	SubsBatch SubMap    //map[string]*Sub 批量数据订阅者通道
	Pubs      PubMap    //map[string]*Pub 发布者通道,消息体为map映射的json
}

func NewRedux() *redux {
	local := NewLocalSub("local_data")
	local.Listen()
	return &redux{
		Local: local,
	}
}

// 添加发布者
//发布者会启动一个协程发布其消息
func (r *redux) AddPub(id string) {
	pub := NewPub()
	r.Pubs.Store(id, pub)
	go r.flowDistribution(id)
}

// 添加订阅
//添加一个订阅者,根据是否为流订阅进行区分储存
func (r *redux) AddSub(so *SubOption) {
	sub := NewSub(so)
	if sub.IsFlow {
		r.SubsFlow.Store(sub.Id, sub)
		defer func() {
			if err := recover(); err != nil {
				r.RemoveSub(so.Id)
				logrus.Error("recover error:", err)
			}
		}()
		go func() {
			sub.SubFlowWork()
		}()
	} else {
		r.SubsBatch.Store(sub.Id, sub)
	}
}

// 删除发布者
func (r *redux) RemovePub(id string) {
	if pub := r.Pubs.Load(id); pub != nil {
		pub.Cancel()
		MsgChanPool.Put(pub.MsgChan)
		r.Pubs.Delete(id)
	}
}

//	删除订阅者
func (r *redux) RemoveSub(id string) {
	if sub := r.SubsFlow.Load(id); sub != nil {
		sub.Cancel()
		MsgStrChanPool.Put(sub.MsgStrChan)
		r.SubsFlow.Delete(id)
	} else if sub := r.SubsBatch.Load(id); sub != nil {
		sub.Cancel()
		MsgStrChanPool.Put(sub.MsgStrChan)
		r.SubsBatch.Delete(id)
	}
}

//将消息分发到流订阅者
func (r *redux) flowDistribution(pubId string) {
	pub := r.Pubs.Load(pubId)
	for {
		select {
		case <-pub.Ctx.Done():
			return
		case msg := <-pub.MsgChan:
			msg["@id"] = pubId
			msgByte, err := json.Marshal(&msg)
			if err != nil {
				logrus.Errorf(err.Error())
			}
			msgStr := string(msgByte)
			r.Local.MsgStrChan <- msgStr
			r.SubsFlow.Range(func(key string, sub *Sub) {
				subChan := sub.MsgStrChan
				subChan <- msgStr
			})
		}
	}
}

//定时分发
//在创建本地订阅者后调用以定时将数据分发给批订阅者
func (r *redux) TimerDistribution(isDrop bool) {
	daily := gron.Every(1 * xtime.Day)
	//weekly := gron.Every(1 * xtime.Week)
	//monthly := gron.Every(30 * xtime.Day)
	//yearly := gron.Every(365 * xtime.Day
	//Timer.AddFunc(gron.Every(tarTime).At(timePoint), func() {
	//
	//})
	Timer := gron.New()
	Timer.AddFunc(daily.At("00:00"), func() {
		go func() {
			for PathKey, fp := range r.Local.GetYesterdayAllFp() {
				r.SubsBatch.Range(func(key string, sub *Sub) {
					defer func() {
						if err := recover(); err != nil {
							r.RemoveSub(sub.Id)
							logrus.Error("recover error:", err)
						}
					}()
					_ = sub.SubBatchWork(PathKey, fp)
				})
				_ = fp.Close()
				if isDrop {
					err := os.Remove(PathKey)
					if err != nil {
						logrus.Errorf("删除%s文件错误:%v", PathKey, err)
					}
				}
				r.Local.DropOneFp(PathKey)
			}
		}()
	})

	Timer.Start()
}
