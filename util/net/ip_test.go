package net

import (
	"testing"

	"github.com/smiecj/go_common/util/log"
	"github.com/stretchr/testify/require"
)

// 测试本地端口是否被占用
func TestGetLocalIP(t *testing.T) {
	// 22: sshd
	ip, err := GetLocalIPV4()
	require.Empty(t, err)
	log.Info("[test] local ip: %s", ip)
}
