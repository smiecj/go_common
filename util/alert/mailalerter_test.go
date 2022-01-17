package alert

import (
	"testing"

	"github.com/smiecj/go_common/util/mail"
)

const (
	testDefaultReceiver = "default_receiver"
	testMailToken       = "token"
	testSender          = "sender"

	testAlertTitle = "alert title"
	testAlertMsg   = "alert msg"
)

func TestSendMail(t *testing.T) {
	sender := mail.NewQQMailSender(mail.MailSenderConf{
		Token:  testMailToken,
		Sender: testSender,
	})
	alerter := GetMailAlerter(sender, testDefaultReceiver)
	alerter.Alert(SetAlertTitleAndMsg(testAlertTitle, testAlertMsg))
}
