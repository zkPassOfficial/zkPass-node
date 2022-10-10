package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterPing(r *gin.RouterGroup) {
	r.GET("/ping", ping)
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"mes": "pong",
	})
}
