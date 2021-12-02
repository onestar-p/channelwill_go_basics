package email

import (
	"channelwill_go_basics/global"
	"crypto/tls"
	"fmt"

	"gopkg.in/gomail.v2"
)

// 邮件发送
type Email struct {
	Host     string // smtp服务器地址
	Port     int    // smtp服务器端口
	Username string // smtp登录的账号
	Fromname string
	Passwd   string // smtp登录的密码，使用生成的授权码
}

func NewEmail() *Email {
	emailConfig := global.ApplicationConfig.EmailInfo
	return &Email{
		Host:     emailConfig.Host,
		Port:     emailConfig.Port,
		Username: emailConfig.Username,
		Fromname: emailConfig.Fromname,
		Passwd:   emailConfig.Passwd,
	}
}

// 发送邮件
func (e *Email) Send(to []string, subject, body string) error {
	if len(to) == 0 {
		return fmt.Errorf("%s", "Need To.")
	}

	m := gomail.NewMessage(
		gomail.SetCharset("utf-8"),
	)
	m.SetHeader("From", m.FormatAddress(e.Username, e.Fromname))
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	d := gomail.NewDialer(e.Host, e.Port, e.Username, e.Passwd) // 设置邮件正文
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	err := d.DialAndSend(m)
	return err
}
