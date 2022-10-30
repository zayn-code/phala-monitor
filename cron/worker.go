package cron

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/shopspring/decimal"
	"log"
	"pha/common"
	"pha/db"
	"pha/global"
	"strconv"
	"strings"
	"time"
)

type worker struct {
	Uuid     string `json:"uuid"`
	Endpoint string `json:"endpoint"`
	Name     string `json:"name"`
	Pid      string `json:"pid"`
	Stake    string `json:"stake"`
}
type workerStates struct {
	Initialized    bool   `json:"initialized"`
	LastMessage    string `json:"lastMessage"`
	MinerAccountId string `json:"minerAccountId"`
	MinerInfoJson  string `json:"minerInfoJson"`
	Status         string `json:"status"`
	Worker         worker `json:"worker"`
}

type workerStatesData struct {
	WorkerStates []workerStates `json:"workerStates"`
}

type workerStatusResp struct {
	Data     workerStatesData `json:"data"`
	HasError bool             `json:"hasError"`
}

type stats struct {
	TotalReward string `json:"totalReward"`
}

type minerInfo struct {
	State string `json:"state"`
	Ve    string `json:"ve"`
	V     string `json:"v"`
	Stats stats  `json:"stats"`
}

func WorkerStart() {
	url := global.PrbConfig.Origin + "/ptp/proxy/" + global.PrbConfig.PeerId + "/GetWorkerStatus"

	//发起请求
	var resp workerStatusResp
	client := resty.New()
	res, err := client.R().
		SetResult(&resp).
		Post(url)
	if err != nil {
		log.Println("request workerStatus error:", err)
	}
	if res.StatusCode() != 200 {
		log.Println("res statusCode error:", res)
	}

	//是否保存数据
	isSave := false
	date := time.Now().Format(global.GoDate)
	var todayData []db.Worker
	global.DB.Where("date = ?", date).Limit(1).Find(&todayData)
	if len(todayData) == 0 {
		yesterday := time.Now().AddDate(0, 0, -1).Format(global.GoDate)
		var yesterdayData db.Worker
		global.DB.Where("date = ?", yesterday).First(&yesterdayData)
		nowYes := time.Now().AddDate(0, 0, -1)
		isSave = yesterdayData.CreatedAt.Before(nowYes)
	}

	//忽略的worker
	var ignoreWorkers WorkerStatusListType
	_ = json.Unmarshal([]byte(db.GetConfig(global.IgnoreWorkers)), &ignoreWorkers)

	alarmContent := ""
	for _, w := range resp.Data.WorkerStates {
		var miner minerInfo
		_ = json.Unmarshal([]byte(w.MinerInfoJson), &miner)

		if isSave {
			saveData(w, date, miner)
		}

		if _, ok := ignoreWorkers[w.Worker.Name]; !ok {
			if miner.State != "MiningIdle" || !strings.Contains(w.LastMessage, "Now the worker should be mining.") {
				if WorkerStatusList[w.Worker.Name]%5 == 0 {
					restartWorker(w)
				}
				if WorkerStatusList[w.Worker.Name]%720 == 0 {
					alarmContent += "<br>\n------<br>\n<b>" + w.Worker.Name + "</b><br>\n" + w.LastMessage + "<br>\n------<br>\n"
				}
				WorkerStatusList[w.Worker.Name]++
			} else {
				WorkerStatusList[w.Worker.Name] = 0
			}
		} else {
			WorkerStatusList[w.Worker.Name] = 0
		}
	}

	if alarmContent != "" {
		common.Alarm(alarmContent)
	}
}

//重启worker
func restartWorker(worker workerStates) {
	log.Println("restart worker " + worker.Worker.Name + ":")
	log.Println(worker.LastMessage)

	type reqDataType struct {
		Ids []string `json:"ids"`
	}

	client := resty.New()
	_, err := client.R().
		SetBody(reqDataType{Ids: []string{worker.Worker.Uuid}}).
		Post(global.PrbConfig.Origin + "/ptp/proxy/" + global.PrbConfig.PeerId + "/RestartWorker")
	if err != nil {
		log.Println("restart worker error:", err)
	}
}

//格式化奖励
func formatReward(s string) float64 {
	strA := strings.Split(s, " ")
	if len(strA) < 2 {
		return 0
	}
	f, _ := strconv.ParseFloat(strA[0], 4)
	if strings.Contains(strA[1], "k") {
		return decimal.NewFromFloat(f).Mul(decimal.NewFromFloat(1000)).InexactFloat64()
	}
	if strings.Contains(strA[1], "m") {
		return decimal.NewFromFloat(f).Mul(decimal.NewFromFloat(1000000)).InexactFloat64()
	}
	return f
}

//保存数据
func saveData(worker workerStates, date string, miner minerInfo) {
	pid, _ := strconv.Atoi(worker.Worker.Pid)
	ve, _ := strconv.ParseFloat(miner.Ve, 8)
	v, _ := strconv.ParseFloat(miner.V, 8)
	totalReward := formatReward(miner.Stats.TotalReward)
	data := db.Worker{
		Uuid:        worker.Worker.Uuid,
		Pid:         uint(pid),
		Name:        worker.Worker.Name,
		Status:      worker.Status,
		State:       miner.State,
		LastMessage: worker.LastMessage,
		Endpoint:    worker.Worker.Endpoint,
		Ve:          ve,
		V:           v,
		TotalReward: totalReward,
		Date:        date,
	}
	if err := global.DB.Create(&data).Error; err != nil {
		log.Println("worker data insert error ", err)
	}
}
