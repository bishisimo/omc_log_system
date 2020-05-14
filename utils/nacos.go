/*
@author '彼时思默'
@time 2020/4/4 20:48
@describe:
*/
package utils

import (
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/utils"
	"github.com/nacos-group/nacos-sdk-go/vo"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"os"
)

type Nacos struct {
	Ip                      string
	Port                    uint
	DataId                  string
	ServerName              string
	ClientConfig            constant.ClientConfig
	ServerConfigs           []constant.ServerConfig
	NamingClient            naming_client.INamingClient
	ConfigClient            config_client.IConfigClient
	RegisterInstanceParam   vo.RegisterInstanceParam
	DeRegisterInstanceParam vo.DeregisterInstanceParam
}

func NewNacos(ip string, port uint, serverName string, nacosIp string) *Nacos {
	logDir := "nacos/logs"
	cacheDir := "nacos/cache"
	if _, err := os.Stat(logDir); !os.IsExist(err) {
		_ = os.MkdirAll(logDir, os.ModePerm)
	}
	if _, err := os.Stat(cacheDir); !os.IsExist(err) {
		_ = os.MkdirAll(cacheDir, os.ModePerm)
	}
	cc := constant.ClientConfig{
		TimeoutMs:      10 * 1000,
		ListenInterval: 30 * 1000,
		BeatInterval:   5 * 1000,
		LogDir:         logDir,
		CacheDir:       cacheDir,
		//NamespaceId:    "public",
	}

	sc := []constant.ServerConfig{
		{
			IpAddr:      nacosIp,
			ContextPath: "/nacos",
			Port:        8848,
		},
	}

	nameClient, err := clients.CreateNamingClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})

	if err != nil {
		logrus.Panic(err)
	}

	cfgClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})
	if err != nil {
		logrus.Panic(err)
	}
	//cluster:="DEFAULT"
	return &Nacos{
		Ip:            ip,
		Port:          port,
		ServerName:    serverName,
		DataId:        serverName + uuid.NewV4().String(),
		ClientConfig:  cc,
		ServerConfigs: sc,
		NamingClient:  nameClient,
		ConfigClient:  cfgClient,
	}
}

/**
注册服务实例
*/
func (n Nacos) Register() bool {
	logrus.Info("开始注册服务实例")
	success, err := n.NamingClient.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          n.Ip,
		Port:        uint64(n.Port),
		ServiceName: n.ServerName,
		Weight:      1,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
	})
	if err != nil {
		logrus.Panic("注册服务实例失败", err)
	}
	logrus.Info("注册服务实例成功")
	return success
}

/**
注销服务实例
*/
func (n Nacos) DeRegister() bool {
	logrus.Info("开始注销服务实例")
	success, err := n.NamingClient.DeregisterInstance(vo.DeregisterInstanceParam{
		Ip:          n.Ip,
		Port:        uint64(n.Port),
		ServiceName: n.ServerName,
		Ephemeral:   true,
	})
	if err != nil {
		logrus.Panic("注销服务实例错误:", err)
	}
	logrus.Info("注销服务实例成功")
	return success
}

/**
服务发现
*/
func (n Nacos) GetService(serverName string) model.Service {
	service, _ := n.NamingClient.GetService(vo.GetServiceParam{
		ServiceName: serverName,
	})
	return service
}

/**
获取所有实例
*/
func (n Nacos) SelectAllInstances(serverName string) []model.Instance {
	instances, err := n.NamingClient.SelectAllInstances(vo.SelectAllInstancesParam{
		ServiceName: serverName,
	})
	if err != nil {
		logrus.Error("获取所有实例失败!")
	}
	return instances
}

/**
获取实例
*/
func (n Nacos) SelectInstances(serverName string) []model.Instance {
	instances, err := n.NamingClient.SelectInstances(vo.SelectInstancesParam{
		ServiceName: serverName,
		HealthyOnly: true,
	})
	if err != nil {
		logrus.Error("获取实例失败!")
	}
	return instances
}

/**
获取一个健康的实例（加权轮训负载均衡）
*/
func (n Nacos) SelectOneHealthyInstance(serverName string) *model.Instance {
	instance, err := n.NamingClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		ServiceName: serverName,
	})
	if err != nil {
		logrus.Error("获取健康实例失败!")
	}
	return instance
}

/**
服务监听
*/
func (n Nacos) Subscribe(serverName string) {
	err := n.NamingClient.Subscribe(&vo.SubscribeParam{
		ServiceName: serverName,
		SubscribeCallback: func(services []model.SubscribeService, err error) {
			logrus.Infof("\n\n callback return services:%s \n\n", utils.ToJsonString(services))
		},
	})
	if err != nil {
		logrus.Error("服务监听失败!", err)
	}
}

/**
取消服务监听
*/
func (n Nacos) Unsubscribe(serverName string) {
	err := n.NamingClient.Unsubscribe(&vo.SubscribeParam{
		ServiceName: serverName,
		SubscribeCallback: func(services []model.SubscribeService, err error) {
			logrus.Infof("\n\n callback return services:%s \n\n", utils.ToJsonString(services))
		},
	})
	if err != nil {
		logrus.Error("取消服务监听失败!", err)
	}
}

/**
发布配置
*/
func (n Nacos) PublishConfig(content string) bool {
	success, err := n.ConfigClient.PublishConfig(vo.ConfigParam{
		DataId:  "dataId",
		Group:   "group",
		Content: content,
	})
	if err != nil {
		logrus.Error("发布配置失败!", err)
	}
	return success
}

/**
删除配置
*/
func (n Nacos) DeleteConfig() bool {
	success, err := n.ConfigClient.DeleteConfig(vo.ConfigParam{
		DataId: "dataId",
		Group:  "group",
	})
	if err != nil {
		logrus.Error("删除配置失败!", err)
	}
	return success
}

/**
获取配置
*/
func (n Nacos) GetConfig() string {
	content, err := n.ConfigClient.GetConfig(vo.ConfigParam{
		DataId: "dataId",
		Group:  "group",
	})
	if err != nil {
		logrus.Error("获取配置失败!", err)
	}
	return content
}

/**
监听配置
*/
func (n Nacos) ListenConfig() {
	err := n.ConfigClient.ListenConfig(vo.ConfigParam{
		DataId: "dataId",
		OnChange: func(namespace, group, dataId, data string) {
			logrus.Infof("dataId:%s\tdata:%s\n", dataId, data)
		},
	})
	if err != nil {
		logrus.Error("监听配置失败!", err)
	}
}
