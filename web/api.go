package web

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"pha/cron"
	"pha/db"
	"pha/global"
)

//保存忽略的worker
func SaveIgnoreWorker(c *gin.Context) {
	var workersMap cron.WorkerStatusListType
	if err := c.BindJSON(&workersMap); err != nil {
		c.JSON(http.StatusBadRequest, workersMap)
		return
	}
	workers, _ := json.Marshal(workersMap)
	err := db.SaveConfig(global.IgnoreWorkers, string(workers))
	if err != nil {
		log.Println("save ignore workers error:", err)
		c.JSON(http.StatusServiceUnavailable, err.Error())
		return
	}
	c.JSON(http.StatusOK, workersMap)
	return
}

//获取workers数据
func GetWorkersData(c *gin.Context) {
	exceptionWorkers := make(cron.WorkerStatusListType)
	for worker, count := range cron.WorkerStatusList {
		if count > 0 {
			exceptionWorkers[worker] = count
		}
	}
	var ignoreWorkers cron.WorkerStatusListType
	_ = json.Unmarshal([]byte(db.GetConfig(global.IgnoreWorkers)), &ignoreWorkers)
	c.JSON(http.StatusOK, map[string]cron.WorkerStatusListType{
		"ignoreWorkers":    ignoreWorkers,
		"exceptionWorkers": exceptionWorkers,
	})
	return
}
