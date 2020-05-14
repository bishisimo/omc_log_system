/*
@author '彼时思默'
@time 2020/4/2 15:16
@describe:
*/
package utils

import (
	"encoding/json"
	"fmt"
	"github.com/bishisimo/storeSystem"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"strings"
	"time"
)

type MsgSt struct {
	StoreSystem store.StoreSystem
	MQC         interface{}
	timePoint   string
	fpMap       map[string]*os.File
	msgChan     chan string
	pathChan    chan string
	Driver      string
}

func NewMsgSt(driver string) *MsgSt {
	var storeSystem store.StoreSystem
	var mqc interface{}
	switch driver {
	case "s3":
		storeSystem = store.NewS3Connector()
	case "cos":
		storeSystem = store.NewCos()
	case "mq":
		fmt.Println("暂未支持!")
	default:
		fmt.Println("暂未支持")
	}
	return &MsgSt{
		StoreSystem: storeSystem,                    // s3c连接器
		MQC:         mqc,                            // 消息队列连接器
		fpMap:       make(map[string]*os.File, 256), //文件句柄池
		pathChan:    make(chan string, 256),         // 文件路径通道
		msgChan:     make(chan string, 10240),       // 消息通道
		Driver:      driver,
	}
}

func (m *MsgSt) Send(msg map[string]interface{}, projectName string) {
	var msgStr string
	var appName string
	appName, _ = msg["soft_name"].(string)
	msgByte, err := json.Marshal(&msg)
	if err != nil {
		logrus.Infof(err.Error())
	}
	msgStr = string(msgByte)
	m.msgChan <- msgStr + "\n"
	// 设置发送对象
	switch m.Driver {
	case "s3":
		m.send2s3(projectName, appName)
	case "s3direct":
		m.send2s3Direct(projectName, appName)
	case "mq":
		fmt.Println("暂未支持...")
	default:
		logrus.Error("不支持的发送对象")
	}
}

func (m *MsgSt) send2s3(projectName string, appName string) {
	localDir := path.Join("out", projectName, appName)
	if _, err := os.Stat(localDir); !os.IsExist(err) {
		_ = os.MkdirAll(localDir, os.ModePerm)
	}
	timePoint := time.Now().Format("2006_01_02")
	if m.timePoint == "" {
		m.timePoint = timePoint
	}
	if m.timePoint != timePoint {
		go func() {
			for updatePath := range m.pathChan {
				descPath := strings.Replace(updatePath, "out", "logs", 1)
				_, _ = m.fpMap[updatePath].Seek(0, 0)
				m.StoreSystem.UploadFileByFP(m.fpMap[updatePath], descPath)
				err := m.fpMap[updatePath].Close()
				if err != nil {
					logrus.Errorf("关闭%s文件错误:%v", updatePath, err)
				}
				err = os.Remove(updatePath)
				if err != nil {
					logrus.Errorf("删除%s文件错误:%v", updatePath, err)
				}
				delete(m.fpMap, updatePath)
			}
		}()
		for key := range m.fpMap {
			m.pathChan <- key
		}
		m.timePoint = timePoint
	}
	localPath := path.Join(localDir, fmt.Sprintf("%s.log", timePoint))

	if m.fpMap[localPath] == nil {
		fp, err := os.OpenFile(localPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
		if err != nil {
			logrus.Errorf("OpenFile %s failure", localPath)
		}
		m.fpMap[localPath] = fp
	}
	ms := <-m.msgChan
	_, err := m.fpMap[localPath].WriteString(ms)
	if err != nil {
		logrus.Errorf("WriteString %s failure: %s", ms, err)
	}
}

func (m *MsgSt) send2s3Direct(projectName string, appName string) {
	descDir := path.Join("logs", projectName, appName)
	ms := <-m.msgChan
	m.StoreSystem.UploadString(ms, descDir)
}
