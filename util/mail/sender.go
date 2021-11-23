// package mail 邮件发送器
package mail

import (
	"fmt"
	"net/smtp"
	"strings"
	"sync"

	"github.com/smiecj/go_common/errorcode"
	"github.com/smiecj/go_common/util/log"
)

var (
	// map
	mailSenderMap = make(map[MailSenderConf]Sender)
	// lock
	mailSenderLock sync.RWMutex
)

// 邮件发送器配置定义
type MailSenderConf struct {
	Host   string
	Port   int
	Token  string
	Sender string
}

// 具体邮件发送配置定义
type mailSendConf struct {
	title       string
	receiverArr []string
	ccArr       []string
	content     string
	nickName    string
}

// 邮件发送结构体定义
type Sender interface {
	Send(...mailSendConfSetter) error
}

// QQ邮箱 发送具体实现
type mailSenderQQImpl struct {
	conf *MailSenderConf
}

// 发送邮件具体配置 设置方法定义
type mailSendConfSetter func(*mailSendConf)

// 设置邮件标题
func SetTitle(title string) func(*mailSendConf) {
	return func(conf *mailSendConf) {
		conf.title = title
	}
}

// 设置邮件内容
func SetContent(content string) func(*mailSendConf) {
	return func(conf *mailSendConf) {
		conf.content = content
	}
}

// 设置邮件接收人
func SetReceiver(receiverArr []string) func(*mailSendConf) {
	return func(conf *mailSendConf) {
		conf.receiverArr = append(conf.receiverArr, receiverArr...)
	}
}

// 添加邮件接收人
func AddReceiver(receiver string) func(*mailSendConf) {
	return func(conf *mailSendConf) {
		conf.receiverArr = append(conf.receiverArr, receiver)
	}
}

// 设置邮件抄送人
func SetCC(ccArr []string) func(*mailSendConf) {
	return func(conf *mailSendConf) {
		conf.ccArr = append(conf.ccArr, ccArr...)
	}
}

// 添加邮件抄送人
func AddCC(cc string) func(*mailSendConf) {
	return func(conf *mailSendConf) {
		conf.ccArr = append(conf.ccArr, cc)
	}
}

// 设置发送人别名
func SetNickName(nickName string) func(*mailSendConf) {
	return func(conf *mailSendConf) {
		conf.nickName = nickName
	}
}

// 发送邮件接口实现
func (impl mailSenderQQImpl) Send(setterArr ...mailSendConfSetter) error {
	conf := new(mailSendConf)
	for _, setter := range setterArr {
		setter(conf)
	}
	// SMTP 密码，不要提交到git 上
	auth := smtp.PlainAuth("", impl.conf.Sender, impl.conf.Token, impl.conf.Host)

	contentType := "Content-Type: text/plain; charset=UTF-8"

	msg := []byte("To: " + strings.Join(conf.receiverArr, ",") + "\r\nFrom: " + conf.nickName +
		"<" + impl.conf.Sender + ">\r\nSubject: " + conf.title + "\r\n" + contentType + "\r\n\r\n" + conf.content)
	err := smtp.SendMail(fmt.Sprintf("%s:%d", impl.conf.Host, impl.conf.Port), auth, impl.conf.Sender, conf.receiverArr, msg)
	if err != nil {
		log.Error("[mailSenderQQImpl.SendMail] send mail failed: %s", err.Error())
		return errorcode.BuildErrorWithMsg(errorcode.SendMailFailed, err.Error())
	}
	log.Info("[mailSenderQQImpl.SendMail] send mail success, sender: %s", impl.conf.Sender)
	return nil
}

// 构建一个 QQ 邮件发送器
func NewQQMailSender(conf MailSenderConf) Sender {
	var sender Sender
	mailSenderLock.RLock()
	sender = mailSenderMap[conf]
	mailSenderLock.RUnlock()

	if nil != sender {
		return sender
	}

	mailSenderLock.Lock()
	defer mailSenderLock.Unlock()

	// mysql 连接能成功创建，并执行 SQL, 才算是创建成功
	impl := new(mailSenderQQImpl)
	if conf.Host == "" {
		conf.Host = "smtp.qq.com"
	}
	if conf.Port == 0 {
		conf.Port = 587
	}
	impl.conf = &conf
	mailSenderMap[conf] = impl
	return impl
}
