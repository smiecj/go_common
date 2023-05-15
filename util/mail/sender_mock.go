package mail

import (
	"github.com/smiecj/go_common/util/log"
)

type mockMailSender struct{}

// 发送邮件，调用 smtp 接口逻辑
func (impl *mockMailSender) Send(setterArr ...mailSendConfSetter) error {
	conf := new(mailSendConf)
	for _, setter := range setterArr {
		setter(conf)
	}

	log.Info("[mockMailSender] send mail: title: %s, msg: %s", conf.title, conf.content)

	return nil
}
