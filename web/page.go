package web

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"pha/db"
	"pha/global"
)

func Workers(c *gin.Context) {
	c.HTML(http.StatusOK, "workers.html", gin.H{})
}

func Income(c *gin.Context) {
	var workerNames []string
	global.DB.Model(&db.Worker{}).Select("name").Find(&workerNames)
	workerNamesStr, _ := json.Marshal(sliceDistinct(workerNames))
	c.HTML(http.StatusOK, "income.html", gin.H{
		"workerNames": string(workerNamesStr),
	})
}

//数组去重
func sliceDistinct(arr []string) []string {
	set := make(map[string]struct{}, len(arr))
	j := 0
	for _, v := range arr {
		_, ok := set[v]
		if ok {
			continue
		}
		set[v] = struct{}{}
		arr[j] = v
		j++
	}

	return arr[:j]
}
