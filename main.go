package main

import (
	"pha/cron"
	"pha/db"
	"pha/web"
)

func main() {
	InitEnv()
	db.InitSqlite()
	cron.InitCron()
	web.InitWeb()
}
