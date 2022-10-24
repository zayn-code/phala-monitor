package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Workers(c *gin.Context) {
	c.HTML(http.StatusOK, "workers.html", gin.H{})
}

func Income(c *gin.Context) {
	c.HTML(http.StatusOK, "income.html", gin.H{})
}
