package alert

import (
	"github.com/smiecj/go_common/errorcode"
	"github.com/smiecj/go_common/util/mail"
)

// mail alerter
type mailAlerter struct {
	mailSender mail.Sender
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
	if len(conf.receiverArr) == 0 {
		return errorcode.BuildError(errorcode.AlertReceiverEmpty)
	}

	return alerter.mailSender.Send(mail.SetReceiver(conf.receiverArr), mail.SetTitle(conf.title), mail.SetContent(conf.msg))
}

// get mail alerter
func GetMailAlerter(sender mail.Sender) Alerter {
	return &mailAlerter{
		mailSender: sender,
	}
}
