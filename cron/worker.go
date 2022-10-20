package cron

import (
	"encoding/json"
	"errors"
	"github.com/go-resty/resty/v2"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"log"
	"pha/common"
	"pha/db"
	"pha/global"
	"strconv"
	"strings"
	"time"
)

const (
	PrbOrigin = "PrbOrigin"
	PeerId    = "PeerId"
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
type WorkerStatusListType map[string]int

var WorkerStatusList WorkerStatusListType
var origin string
var peerId string

func init() {
	origin = common.GetEnvDefault(PrbOrigin, "http://192.168.2.239:3000")
	peerId = common.GetEnvDefault(PeerId, "")
	if peerId == "" {
		common.ErrorExit(common.WithErrorMsg("peerId is required"))
	}
	WorkerStatusList = make(map[string]int)
}

func WorkerStart() {
	url := origin + "/ptp/proxy/" + peerId + "/GetWorkerStatus"

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
	if firstErr := global.DB.Where("date = ?", date).First(&db.Worker{}).Error; errors.Is(firstErr, gorm.ErrRecordNotFound) {
		yesterday := time.Now().AddDate(0, 0, -1).Format(global.GoDate)
		var yesterdayData db.Worker
		global.DB.Where("date = ?", yesterday).First(&yesterdayData)
		nowYes := time.Now().AddDate(0, 0, -1)
		isSave = yesterdayData.CreatedAt.Before(nowYes)
	}

	//忽略的worker
	var ignoreWorkers map[string]int
	_ = json.Unmarshal([]byte(db.GetConfig(global.IgnoreWorkers)), &ignoreWorkers)

	for _, w := range resp.Data.WorkerStates {
		if isSave {
			saveData(w, date)
		}

		if _, ok := ignoreWorkers[w.Worker.Name]; !ok {
			if w.Status != "S_MINING" || !strings.Contains(w.LastMessage, "Now the worker should be mining.") {
				if WorkerStatusList[w.Worker.Name]%5 == 0 {
					restartWorker(w.Worker.Uuid, w.LastMessage)
				}
				if WorkerStatusList[w.Worker.Name]%720 == 0 {
					common.Alarm(w.Worker.Name + "\n" + w.LastMessage + "\n\n")
				}
				WorkerStatusList[w.Worker.Name]++
			} else {
				WorkerStatusList[w.Worker.Name] = 0
			}
		}
	}

}

//重启worker
func restartWorker(uuid string, msg string) {
	log.Println("restart worker:")
	log.Println(msg)

	type reqDataType struct {
		Ids []string `json:"ids"`
	}

	client := resty.New()
	_, err := client.R().
		SetBody(reqDataType{Ids: []string{uuid}}).
		Post(origin + "/ptp/proxy/" + peerId + "/RestartWorker")
	if err != nil {
		log.Println("restart worker error:", err)
	}
}

//格式化奖励
func formatReward(s string) float64 {
	strA := strings.Split(s, " ")
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
func saveData(worker workerStates, date string) {
	var miner minerInfo
	_ = json.Unmarshal([]byte(worker.MinerInfoJson), &miner)

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
