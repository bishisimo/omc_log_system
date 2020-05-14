/*
@author '彼时思默'
@time 2020/5/11 上午10:22
@describe:订阅者
*/
package core

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bishisimo/store"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"
)

var SubSupport = map[string][]string{
	"mq":   {"flow", "batch"},
	"s3":   {"flow", "batch"},
	"cos":  {"flow", "batch"},
	"http": {"flow", "batch"},
	"std":  {"flow"},
}

//订阅者
type Sub struct {
	*SubOption
	Ctx        context.Context
	Cancel     context.CancelFunc
	MsgStrChan chan string
}

func NewSub(so *SubOption) *Sub {
	ctx, cancel := context.WithCancel(context.Background())
	return &Sub{
		SubOption:  so,
		MsgStrChan: MsgStrChanPool.Get(),
		Ctx:        ctx,
		Cancel:     cancel,
	}
}

//消息队列流订阅处理
func (s *Sub) MqFlowWork() {
	for {
		msg, ok := <-s.MsgStrChan
		if !ok {
			return
		}
		_ = msg
	}
	//TODO:实现mq的流传输
}

//对象储存流订阅处理
func (s *Sub) OsFlowWork() {
	var sf store.StoreFace
	switch s.AcceptType {
	case "s3":
		sf = S3Pool.Get()
	case "cos":
		sf = CosPool.Get()
	default:
		return
	}
	for {
		select {
		case <-s.Ctx.Done():
			switch s.AcceptType {
			case "s3":
				S3Pool.Put(sf)
			case "cos":
				CosPool.Put(sf)
			}
			return
		case msgStr := <-s.MsgStrChan:
			var msg map[string]interface{}
			err := json.Unmarshal([]byte(msgStr), &msg)
			if err != nil {
				logrus.Error("json unmarshal error:", err)
			}
			today := time.Now().Format("2006-01-02")
			dir := path.Join("logs", msg["@id"].(string), msg["soft_name"].(string), today)
			var fileName string
			if fileName = msg["atom_id"].(string); fileName == "" {
				fileName = "default"
			}
			descPath := path.Join(dir, fileName)
			sf.UploadString(msgStr, descPath)
		}
	}
}

//HTTP流订阅处理
func (s *Sub) HttpFlowWork() {
	/**
	http处理数据
	*/
	for {
		msg, ok := <-s.MsgStrChan
		if !ok {
			return
		}
		request := RequestPool.Get()
		defer RequestPool.Put(request)
		res, err := request.SetBody(msg).Post(s.Url)
		if err != nil {
			logrus.Error("request error:", err)
		}
		logrus.Info(s.Url, res.Status())
	}

}
func (s *Sub) StdFlowWork() {
	for {
		msg, ok := <-s.MsgStrChan
		if !ok {
			return
		}
		fmt.Println("msg:", msg)
	}
}

func (s *Sub) MqBatchWork(localPath string, fp *os.File) {
	//TODO:实现mq的文件
}
func (s *Sub) OsBatchWork(localPath string, fp *os.File) {
	var sf store.StoreFace
	switch s.AcceptType {
	case "s3":
		sf = store.NewS3Connector(s.AccessOption)
	case "cos":
		sf = store.NewCosConnector(s.AccessOption)
	default:
		return
	}
	descPath := strings.Replace(localPath, "local_data", "logs", 1)
	_, _ = fp.Seek(0, 0)
	sf.UploadFileByFP(fp, descPath)
}

//http处理文件数据
func (s *Sub) HttpBatchWork(localPath string, fp *os.File) {
	_, _ = fp.Seek(0, 0)
	body, err := ioutil.ReadAll(fp)
	if err != nil {
		logrus.Error("read error:", err)
	}
	request := RequestPool.Get()
	defer RequestPool.Put(request)
	res, err := request.SetBody(body).Post(s.Url)
	if err != nil {
		logrus.Error("request error:", err)
	}
	logrus.Info(s.Url, res.Status())
}

func (s *Sub) SubFlowWork() {
	switch s.AcceptType {
	case "mq":
		s.MqFlowWork()
	case "s3", "cos":
		s.OsFlowWork()
	case "http":
		s.HttpFlowWork()
	case "std":
		s.StdFlowWork()
	default:
		return
	}
}
func (s *Sub) SubBatchWork(localPath string, fp *os.File) error {
	switch s.AcceptType {
	case "mq":
		s.MqBatchWork(localPath, fp)
		return nil
	case "s3", "cos":
		s.OsBatchWork(localPath, fp)
		return nil
	case "http":
		s.HttpBatchWork(localPath, fp)
		return nil
	default:
		return fmt.Errorf(s.AcceptType)
	}
}
