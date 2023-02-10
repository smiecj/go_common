package alert

import (
	"fmt"

	"github.com/smiecj/go_common/errorcode"
	"github.com/smiecj/go_common/util/mail"
)

// mail alerter
type mailAlerter struct {
	commonTitle string
	mailSender  mail.Sender
}

// set common title
func (alerter *mailAlerter) SetCommonTitle(title string) {
	alerter.commonTitle = title
}

// send an email alert
func (alerter *mailAlerter) Alert(alertConfSetterArr ...alertConfSetter) error {
	conf := new(alertConf)
	for _, setter := range alertConfSetterArr {
		setter(conf)
	}

	if conf.msg == "" {
		return errorcode.BuildError(errorcode.AlertMsgEmpty)
	}

	title := conf.title
	if alerter.commonTitle != "" {
		title = fmt.Sprintf("%s %s", alerter.commonTitle, title)
	}
	return alerter.mailSender.Send(mail.SetReceiver(conf.receiverArr), mail.SetTitle(title), mail.SetContent(conf.msg))
}

// get mail alerter
func GetMailAlerter(sender mail.Sender) Alerter {
	return &mailAlerter{
		mailSender: sender,
	}
}
