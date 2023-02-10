// package mail 邮件发送器
package mail

import (
	"fmt"
	"net/smtp"
	"strings"
	"sync"

	"github.com/smiecj/go_common/config"
	"github.com/smiecj/go_common/errorcode"
	"github.com/smiecj/go_common/util/log"
)

const (
	// 邮件发送器初始化配置中，收件人中的默认分隔符（逗号）
	receiverSplitor = ","
	// 配置中心读取邮件发送配置，默认使用的space
	// 可能会造成多个项目误共用配置的问题，在项目使用中，尽量通过配置文件路径区分不同的项目配置
	mailDefaultConfigSpace = "mail"
)

var (
	// map
	mailSenderMap = make(map[mailSenderInitConf]Sender)
	// lock
	mailSenderLock sync.RWMutex
)

// 邮件发送器配置定义
type mailSenderInitConf struct {
	Host          string `yaml:"host"`
	Port          int    `yaml:"port"`
	Token         string `yaml:"token"`
	Sender        string `yaml:"sender"`
	Receiver      string `yaml:"receiver"`
	SendSeparatly bool   `yaml:"send_separatly"`
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

// SMTP邮箱 发送具体实现
type mailSenderSMTPImpl struct {
	conf mailSenderInitConf
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
func (impl mailSenderSMTPImpl) Send(setterArr ...mailSendConfSetter) error {
	conf := new(mailSendConf)
	for _, setter := range setterArr {
		setter(conf)
	}
	// default receiver use sender init conf
	if len(conf.receiverArr) == 0 {
		conf.receiverArr = strings.Split(strings.TrimSpace(impl.conf.Receiver), receiverSplitor)
	}

	// 设定在发送给多个收件人的时候分开发送
	// 避免不相关的收件人都看到信息
	// 后续: 可以考虑通过 hook 钩子的方式，配置不同的发送通道
	if impl.conf.SendSeparatly {
		for _, currentReceiver := range conf.receiverArr {
			err := impl.send(conf, []string{currentReceiver})
			if nil != err {
				return err
			}
		}
	} else {
		return impl.send(conf, conf.receiverArr)
	}

	return nil
}

// 发送邮件，调用 smtp 接口逻辑
func (impl mailSenderSMTPImpl) send(conf *mailSendConf, receiverArr []string) error {
	auth := smtp.PlainAuth("", impl.conf.Sender, impl.conf.Token, impl.conf.Host)
	contentType := "Content-Type: text/plain; charset=UTF-8"
	msg := []byte("To: " + strings.Join(receiverArr, ",") + "\r\nFrom: " + conf.nickName +
		"<" + impl.conf.Sender + ">\r\nSubject: " + conf.title + "\r\n" + contentType + "\r\n\r\n" + conf.content)
	err := smtp.SendMail(fmt.Sprintf("%s:%d", impl.conf.Host, impl.conf.Port), auth, impl.conf.Sender, receiverArr, msg)
	if err != nil {
		log.Error("[mailSenderSMTPImpl.SendMail] send mail failed: %s", err.Error())
		return errorcode.BuildErrorWithMsg(errorcode.SendMailFailed, err.Error())
	}
	log.Info("[mailSenderSMTPImpl.SendMail] send mail success, sender: %s, receiver: %v", impl.conf.Sender, receiverArr)
	return nil
}

// 构建一个 SMTP 邮件发送器
func NewSMTPMailSender(configManager config.Manager) (Sender, error) {
	var sender Sender
	mailSenderLock.RLock()
	senderInitConf := mailSenderInitConf{}
	configManager.Unmarshal(mailDefaultConfigSpace, &senderInitConf)
	sender = mailSenderMap[senderInitConf]
	mailSenderLock.RUnlock()

	if nil != sender {
		return sender, nil
	}

	mailSenderLock.Lock()
	defer mailSenderLock.Unlock()

	impl := new(mailSenderSMTPImpl)
	if senderInitConf.Host == "" {
		return nil, errorcode.BuildErrorWithMsg(errorcode.SenderInitFailed, "smtp host is empty")
	}
	if senderInitConf.Port == 0 {
		senderInitConf.Port = 587
	}

	// token 和 发送者 不能为空
	if senderInitConf.Token == "" {
		return nil, errorcode.BuildErrorWithMsg(errorcode.SenderInitFailed, "token is empty")
	}
	if senderInitConf.Sender == "" {
		return nil, errorcode.BuildErrorWithMsg(errorcode.SenderInitFailed, "sender is empty")
	}

	impl.conf = senderInitConf
	mailSenderMap[senderInitConf] = impl
	return impl, nil
}
