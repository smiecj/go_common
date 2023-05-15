package mail

import (
	"testing"

	yamlconfig "github.com/smiecj/go_common/config/yaml"
	"github.com/smiecj/go_common/util/file"
	"github.com/stretchr/testify/require"
)

const (
	localConfigFile = "conf_local.yaml"
)

// 测试 通过 SMTPMailSender 发送一封邮件
func TestSendMail(t *testing.T) {
	configManager, err := yamlconfig.GetYamlConfigManager(file.FindFilePath(localConfigFile))
	require.Empty(t, err)

	sender, err := NewSMTPMailSender(configManager)
	require.Empty(t, err)

	err = sender.Send(SetTitle("test_title"), SetContent("test_content"), SetNickName("smiecj"))
	require.Equal(t, nil, err)
}
