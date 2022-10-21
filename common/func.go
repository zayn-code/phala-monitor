package common

import (
	"gopkg.in/gomail.v2"
	"log"
	"os"
	"pha/global"
)

// 获取环境变量信息
func GetEnvDefault(key, defVal string) string {
	val, ex := os.LookupEnv(key)
	if !ex {
		return defVal
	}
	return val
}

//报警
func Alarm(content string) {
	log.Println("start alarm:", content)
	SendMail(content)
}

//发送邮件
func SendMail(body string) {
	m := gomail.NewMessage()
	m.SetHeader("From", global.MailConfig.From)
	m.SetHeader("To", global.MailConfig.To...)
	m.SetHeader("Subject", global.MailConfig.Subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(global.MailConfig.Host, global.MailConfig.Port, global.MailConfig.Username, global.MailConfig.Password)

	if err := d.DialAndSend(m); err != nil {
		log.Println("send mail error:", err)
	}
}
