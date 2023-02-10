// package alert 告警专用，包含告警发送器, 后续迁移到公共仓库中
package alert

// 告警发送器接口定义
// 可基于此接口 实现邮件告警、企业微信告警等功能
type Alerter interface {
	Alert(...alertConfSetter) error
	SetCommonTitle(string)
}

// 告警发送参数定义
// 可设置告警接收人、标题和消息
type alertConf struct {
	title       string
	msg         string
	receiverArr []string
}

// 告警发送配置设置方法
type alertConfSetter func(*alertConf)

// 设置发送标题和消息
func SetAlertTitleAndMsg(title, msg string) alertConfSetter {
	return func(conf *alertConf) {
		conf.title, conf.msg = title, msg
	}
}

// 设置告警消息
func SetAlertMsg(msg string) alertConfSetter {
	return func(conf *alertConf) {
		conf.msg = msg
	}
}

// 设置告警接收人 (单个)
func SetReceiver(receiver string) alertConfSetter {
	return func(conf *alertConf) {
		conf.receiverArr = []string{receiver}
	}
}

// 设置告警接收人 (多个)
func SetReceiverArr(receiverArr []string) alertConfSetter {
	return func(conf *alertConf) {
		conf.receiverArr = make([]string, 0)
		conf.receiverArr = append(conf.receiverArr, receiverArr...)
	}
}
