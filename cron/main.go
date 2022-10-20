package cron

import (
	"github.com/robfig/cron"
	"log"
)

func InitCron() {
	c := cron.New()
	err := c.AddFunc("0 */1 * * * *", WorkerStart)
	if err != nil {
		log.Println("worker error:", err)
	}
	c.Start()
}
