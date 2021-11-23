package mail

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	testMailSender   = "xxx@qq.com"
	testMailReceiver = "xxx@qq.com"
	testMailToken    = "xxx"
)

// 测试 发送一封 QQ 邮件 token 在提前代码前要屏蔽
func TestSendMail(t *testing.T) {
	sender := NewQQMailSender(MailSenderConf{
		Token:  testMailToken,
		Sender: testMailSender,
	})
	err := sender.Send(AddReceiver(testMailReceiver), SetTitle("test_title"), SetContent("test_content"), SetNickName("smiecj"))
	require.Equal(t, nil, err)
}
