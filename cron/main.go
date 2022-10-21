package cron

import (
	"github.com/robfig/cron"
	"log"
)

var WorkerStatusList WorkerStatusListType

func init() {
	WorkerStatusList = make(map[string]int)
}

func InitCron() {
	c := cron.New()
	err := c.AddFunc("0 */1 * * * *", WorkerStart)
	if err != nil {
		log.Println("worker error:", err)
	}
	c.Start()
}
