package web

import (
	"github.com/gin-gonic/gin"
)

func InitPageRouter(pageRouter *gin.RouterGroup) {
	pageRouter.Use()
	{
		pageRouter.GET("/workers", Workers)
	}
}

func InitApiRouter(apiRouter *gin.RouterGroup) {
	apiRouter.Use()
	{
		apiRouter.POST("/saveIgnoreWorkers", SaveIgnoreWorker)
		apiRouter.GET("/getWorkersData", GetWorkersData)
	}
}
