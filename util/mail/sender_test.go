package mail

import (
	"testing"

	"github.com/smiecj/go_common/config"
	"github.com/stretchr/testify/require"
)

// 测试 通过 QQMailSender 发送一封邮件
func TestSendMail(t *testing.T) {
	configManager, err := config.GetYamlConfigManager("/tmp/mailconf.yml")
	require.Empty(t, err)

	sender, err := NewQQMailSender(configManager)
	require.Empty(t, err)

	err = sender.Send(SetTitle("test_title"), SetContent("test_content"), SetNickName("smiecj"))
	require.Equal(t, nil, err)
}
