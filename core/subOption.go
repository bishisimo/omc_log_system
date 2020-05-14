/*
@author '彼时思默'
@time 2020/4/29 上午11:32
@describe:订阅者配置数据
*/
package core

import "github.com/bishisimo/store"

var SupportSubType = map[string]bool{
	"mq":   true,
	"s3":   true,
	"cos":  true,
	"http": true,
}

func NewSubOption() *SubOption {
	return &SubOption{
		Id:         "",
		AcceptType: "",
		IsFlow:     false,
		MqOption:   MqOption{},
		OsOption:   OsOption{},
		HttpOption: HttpOption{},
	}
}

type SubOption struct {
	Id         string `json:"id"`
	AcceptType string `json:"acceptType"`
	IsFlow     bool   `json:"isFlow"`
	MqOption
	OsOption
	HttpOption
	StdOption
}

func (s SubOption) GetInfo() interface{} {
	switch s.AcceptType {
	case "mq":
		return s.MqOption
	case "s3", "cos":
		return s.OsOption
	case "http":
		return s.HttpOption
	default:
		return ""
	}
}

//配置完整性检查
func (o *SubOption) IsInfoRight() bool {
	if o.Id == "" || o.AcceptType == "" {
		return false
	}
	switch o.AcceptType {
	case "std":
		return true
	case "mq":
		return o.MqOption.IsInfoRight()
	case "s3", "cos":
		return o.OsOption.IsInfoRight()
	case "http":
		return o.HttpOption.IsInfoRight()
	default:
		return false
	}
}

type MqOption struct {
	Hosts []string `json:"hosts"`
	Port  int      `json:"port"`
	Topic string   `json:"topic"`
}

//消息队列配置完整性检查
func (o *MqOption) IsInfoRight() bool {
	if o.Hosts == nil || len(o.Hosts) == 0 || o.Port == 0 || o.Topic == "" {
		return false
	}
	return true
}

type OsOption struct {
	store.AccessOption
}

//对象储存配置完整性检查
func (o *OsOption) IsInfoRight() bool {
	if o.Id == "" || o.Secret == "" || o.Endpoint == "" || o.Region == "" || o.Bucket == "" {
		return false
	}
	return true
}

//	http配置
type HttpOption struct {
	Url   string `json:"path"`
	Field string `json:"field"`
}

//http配置完整性检查
func (o *HttpOption) IsInfoRight() bool {
	if o.Url == "" {
		return false
	}
	if o.Field == "" {
		o.Field = "data"
	}
	return true
}

type StdOption struct {
}
