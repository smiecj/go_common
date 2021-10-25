package mail

import (
	"fmt"
	"net/smtp"
	"strings"
)

// todo: 设计 邮件配置文件 和 邮件配置结构体定义 ，后续做成通用的框架配置

/**
 * sender: 发件人邮箱
 * mailPassword: 发件人的SMTP密码
 * subject: 邮件标题
 * sendContent: 发送邮件内容
 */
func SendMail(sender string, mailPassword string, nickname string, subject string,
	sendContent string, receiver []string) {
	smtpHost := "smtp.qq.com"
	smtpPort := ":587"
	// SMTP 密码，不要提交到git 上
	auth := smtp.PlainAuth("", sender, mailPassword, smtpHost)
	user := sender

	contentType := "Content-Type: text/plain; charset=UTF-8"

	msg := []byte("To: " + strings.Join(receiver, ",") + "\r\nFrom: " + nickname +
		"<" + user + ">\r\nSubject: " + subject + "\r\n" + contentType + "\r\n\r\n" + sendContent)
	err := smtp.SendMail(smtpHost+smtpPort, auth, user, receiver, msg)
	if err != nil {
		fmt.Printf("发送邮件失败！错误信息: %v\n", err)
	}
}
