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

type errorInfo struct {
	Msg  string
	Code int
}

//异常状态码赋值
func WithErrorCode(code int) func(info *errorInfo) {
	return func(info *errorInfo) {
		info.Code = code
	}
}

//异常信息赋值
func WithErrorMsg(msg string) func(info *errorInfo) {
	return func(info *errorInfo) {
		info.Msg = msg
	}
}

//异常退出程序
func ErrorExit(withs ...func(info *errorInfo)) {
	info := errorInfo{
		Msg:  "",
		Code: 1,
	}
	for _, with := range withs {
		with(&info)
	}
	os.Exit(info.Code)
}

func Alarm(content string) {
	log.Println("start alarm:", content)
}
