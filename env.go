package main

import (
	"pha/common"
	"pha/global"
	"strconv"
	"strings"
)

//初始化env的配置
func InitEnv() {
	//初始化Prb配置
	global.PrbConfig.Origin = common.GetEnvDefault("PRB_ORIGIN", "http://127.0.0.1:3000")
	global.PrbConfig.PeerId = common.GetEnvDefault("PRB_PEER_ID", "")
	if global.PrbConfig.PeerId == "" {
		panic("peerId is required")
	}

	//初始化邮箱配置
	global.MailConfig.From = common.GetEnvDefault("MAIL_FROM", "from@example.com")
	global.MailConfig.To = strings.Split(common.GetEnvDefault("MAIL_TO", "to@example.com"), ",")
	global.MailConfig.Subject = common.GetEnvDefault("MAIL_SUBJECT", "phala alarm")
	global.MailConfig.Host = common.GetEnvDefault("MAIL_HOST", "smtp.163.com")
	port, err := strconv.Atoi(common.GetEnvDefault("MAIL_PORT", "465"))
	if err != nil {
		panic("mail port wrong!")
	}
	global.MailConfig.Port = port
	global.MailConfig.Username = common.GetEnvDefault("MAIL_USERNAME", "from@example.com")
	//邮箱平台的授权码
	global.MailConfig.Password = common.GetEnvDefault("MAIL_PASSWORD", "")
}
