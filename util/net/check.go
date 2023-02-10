package net

import (
	"net"
	"strconv"
	"time"
)

// 检查本机指定端口是否被占用
func CheckLocalPortIsUsed(port int) bool {
	portStr := strconv.Itoa(port)
	address := net.JoinHostPort("127.0.0.1", portStr)
	conn, err := net.DialTimeout("tcp", address, time.Second)
	if err != nil {
		// 连接失败，暂时定为未被占用
		return false
	} else {
		defer conn.Close()
		if conn != nil {
			return true
		} else {
			return false
		}
	}
}
