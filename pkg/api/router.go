package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Default(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"mes": "zk-node version 0.1",
	})
}

func Register(g *gin.Engine) {

	g.GET("/", Default)
	g.POST("/", Default)

	v1 := g.Group("/api/v1")

	RegisterPing(v1)
	// RegisterAuth(v1)

}
