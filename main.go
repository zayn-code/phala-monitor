package main

import (
	"pha/cron"
	"pha/db"
	"pha/web"
)

func main() {
	db.InitSqlite()
	cron.InitCron()
	web.InitWeb()
}
