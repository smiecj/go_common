package alert

import (
	"testing"

	"github.com/smiecj/go_common/util/mail"
	"github.com/stretchr/testify/require"
)

const (
	testMailToken = "token"
	testSender    = "sender"
	testReceiver  = "receiver"

	testAlertTitle = "alert title"
	testAlertMsg   = "alert msg"
)

func TestSendMail(t *testing.T) {
	sender := mail.NewQQMailSender(mail.MailSenderConf{
		Token:  testMailToken,
		Sender: testSender,
	})
	alerter := GetMailAlerter(sender)
	err := alerter.Alert(SetAlertTitleAndMsg(testAlertTitle, testAlertMsg), SetReceiver(testReceiver))
	require.Empty(t, err)
}
