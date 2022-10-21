package common

import (
	"log"
	"os"
)

// 获取环境变量信息
func GetEnvDefault(key, defVal string) string {
	val, ex := os.LookupEnv(key)
	if !ex {
		return defVal
	}
	return val
}

func Alarm(content string) {
	log.Println("start alarm:", content)
}
