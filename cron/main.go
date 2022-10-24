package cron

import (
	"github.com/robfig/cron"
	"log"
)

type WorkerStatusListType map[string]int

var WorkerStatusList WorkerStatusListType

func init() {
	WorkerStatusList = make(WorkerStatusListType)
}

func InitCron() {
	c := cron.New()
	err := c.AddFunc("0 */1 * * * *", WorkerStart)
	if err != nil {
		log.Println("worker error:", err)
	}
	c.Start()
}
