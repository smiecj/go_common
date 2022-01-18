// package net 和网络相关的工具方法
package net

import (
	"net"

	"github.com/smiecj/go_common/errorcode"
)

// 获取本机 ipv4 地址
func GetLocalIPV4() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", errorcode.BuildErrorWithMsg(errorcode.NetLocalIPGetFailed, err.Error())
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", errorcode.BuildErrorWithMsg(errorcode.NetLocalIPGetFailed, err.Error())
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue
			}
			return ip.String(), nil
		}
	}
	return "", errorcode.BuildError(errorcode.NetLocalIPGetFailed)
}
