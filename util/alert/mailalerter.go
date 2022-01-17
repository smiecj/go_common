package alert

import (
	"github.com/smiecj/go_common/util/mail"
)

// mail alerter
type mailAlerter struct {
	mailSender      mail.Sender
	defaultReceiver string
}

// send an email alert
func (alerter *mailAlerter) Alert(alertConfSetterArr ...alertConfSetter) error {
	conf := new(alertConf)
	for _, setter := range alertConfSetterArr {
		setter(conf)
	}

	if conf.msg == "" {
		// return errorcode.BuildError()
		return nil
	}

	// if receiver is not set, use default receiver
	var err error
	if len(conf.receiverArr) > 0 {
		err = alerter.mailSender.Send(mail.SetReceiver(conf.receiverArr), mail.SetTitle(conf.title), mail.SetContent(conf.msg))
	} else {
		err = alerter.mailSender.Send(mail.AddReceiver(alerter.defaultReceiver), mail.SetTitle(conf.title), mail.SetContent(conf.msg))
	}

	return err
}

// get mail alerter
func GetMailAlerter(sender mail.Sender, defaultReceiver string) Alerter {
	return &mailAlerter{
		mailSender:      sender,
		defaultReceiver: defaultReceiver,
	}
}
