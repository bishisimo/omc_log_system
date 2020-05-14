/*
@author '彼时思默'
@time 2020/3/24 15:41
@describe:
*/
package main

import (
	"fmt"
	"github.com/bishisimo/omc_log_system/proxy"
	"os"
	"os/signal"
	"syscall"
)

var exitChan = make(chan string)

func exitListen() {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	exitChan <- fmt.Sprintf("捕捉到信号:%v", <-sigChan)
}
func main() {
	//instanceIp := flag.String("ip", "automatic acquisition", "服务实例的ip地址")
	//port := flag.Uint("port", 0, "服务实例的端口")
	//serverName := flag.String("sn", "omcg-logs", "服务名")
	//nacosIp := flag.String("nip", "", "nacos的ip地址")
	//isRetainLocal := flag.Bool("retain", false, "是否保留本地文件")
	//flag.Parse()
	//retainLocal:=""
	//if *isRetainLocal{
	//	retainLocal="true"
	//}
	//_=os.Setenv("IsRetainLocal",retainLocal)
	//
	//if *instanceIp == "automatic acquisition" {
	//	netIp := utils.NewNetIp()
	//	instanceIp = &netIp.Intranet
	//	logrus.Infof("未指定ip,使用默认ip:%v", *instanceIp)
	//}
	//if *port == 0 {
	//	logrus.Panic("未指定端口,实例创建失败!")
	//}
	//if *nacosIp == "" {
	//	logrus.Panic("未指定nacos的IP,实例创建失败!")
	//}
	//go exitListen()
	//nacos := utils.NewNacos(*instanceIp, *port, *serverName, *nacosIp)
	//nacos.Register()
	//defer func() {
	//	if err:=recover();err!=nil{
	//		nacos.DeRegister()
	//	}
	//}()
	//_ = agent.Listen(agent.Options{})
	//go proxy.CreateProxy(8080)
	//logrus.Infof("服务中止,退出!%v", <-exitChan)
	proxy.CreateProxy(8080)
}
