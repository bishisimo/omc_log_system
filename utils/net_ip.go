/*
@author '彼时思默'
@time 2020/4/5 18:47
@describe:
*/
package utils

import (
	"github.com/sirupsen/logrus"
	"net"
	"strings"
)

type NetIp struct {
	Intranet string
	Extranet string
	AllIp    map[string]string
}

func NewNetIp() *NetIp {
	validNet := map[string]bool{
		"WLAN":     true,
		"以太网":      true,
		"wlan0":    true,
		"eth0":     true,
		"eth1":     true,
		"es33":     true,
		"enp4s0f1": true,
		"wlp3s0":   true,
		"enp1s0":   true,
	}
	var intranet string
	var extranet string
	allIp := make(map[string]string)
	netInterfaces, err := net.Interfaces()
	if err != nil {
		logrus.Error("net.Interfaces failed, err:", err.Error())
	}
	for i := 0; i < len(netInterfaces); i++ {
		if netInterfaces[i].Flags&net.FlagUp != 0 {
			keys := strings.Split(netInterfaces[i].Name, " ")
			key := keys[len(keys)-1]
			addrs, _ := netInterfaces[i].Addrs()
			for _, address := range addrs {
				if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
					if ipNet.IP.To4() != nil && validNet[netInterfaces[i].Name] {
						intranet = ipNet.IP.String()
						allIp[key] = intranet
					} else {
						allIp[key] = ipNet.IP.String()
					}
				}
			}
		}
	}

	//resp, err := http.Get("http://localhost/api/ip/raw")
	//if err != nil {
	//	logrus.Error(err.Error())
	//}
	//if resp != nil {
	//	defer resp.Subs.Close()
	//	content, _ := ioutil.ReadAll(resp.Subs)
	//	extranet = string(content)
	//	allIp["ex"] = extranet
	//}
	return &NetIp{
		Intranet: intranet,
		Extranet: extranet,
		AllIp:    allIp,
	}
}
