package web

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"log"
	"net/http"
	"pha/cron"
	"pha/db"
	"pha/global"
	"time"
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

type DayIncome map[string]float64

//获取worker收益
func GetWorkerIncome(c *gin.Context) {
	type incomeParams struct {
		WorkerNames []string `json:"workerNames"`
	}
	var params incomeParams
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	if len(params.WorkerNames) == 0 {
		c.JSON(http.StatusBadRequest, "worker名称不能为空！")
		return
	}
	data := make(map[string]DayIncome)
	for _, workerName := range params.WorkerNames {
		var workerDays []db.Worker
		global.DB.Where("name = ?", workerName).Order("date ASC").Find(&workerDays)
		data[workerName] = getDayIncome(workerDays)
	}

	c.JSON(http.StatusOK, data)
	return
}

//获取每日收益，不足一天的不计算
func getDayIncome(workerDays []db.Worker) DayIncome {
	dayIncome := make(DayIncome)
	for _, workerDay := range workerDays {
		dayIncome[workerDay.Date] = workerDay.TotalReward
	}
	newDayIncome := make(DayIncome)
	for day, income := range dayIncome {
		dayTime, err := time.Parse(global.GoDate, day)
		if err != nil {
			log.Println("date convert time error：", err)
		}
		addOneDay := dayTime.AddDate(0, 0, 1)
		if addOneIncome, ok := dayIncome[addOneDay.Format(global.GoDate)]; ok {
			newDayIncome[day] = decimal.NewFromFloat(addOneIncome).Sub(decimal.NewFromFloat(income)).InexactFloat64()
		}
	}
	return newDayIncome
}
