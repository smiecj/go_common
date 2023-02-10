package net

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// 测试本地端口是否被占用
func TestIsLocalPortInUsed(t *testing.T) {
	// 22: sshd
	isUsed := CheckLocalPortIsUsed(22)
	require.True(t, isUsed)
}
