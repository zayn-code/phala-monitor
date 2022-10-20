package web

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"pha/common"
)

func InitWeb() {
	r := gin.Default()

	r.StaticFS("static", http.Dir("web/static"))
	r.Delims("{[", "]}")
	r.LoadHTMLGlob("web/page/*")

	r.GET("/", func(context *gin.Context) {
		context.Redirect(http.StatusFound, "/page/workers")
	})

	pageRouter := r.Group("/page")
	InitPageRouter(pageRouter)

	apiRouter := r.Group("/api")
	InitApiRouter(apiRouter)

	err := r.Run(":8080")
	if err != nil {
		log.Println("init web error:", err)
		common.ErrorExit()
	}
}
