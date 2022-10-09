package api

import (
	"net/http"
	"zkpass-node/pkg/session"

	"github.com/gin-gonic/gin"
)

func Default(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"mes": "zk-node version 0.1",
	})
}

func Route(sm *session.Manager, g *gin.Engine) {

	g.GET("/", Default)
	g.POST("/", Default)
	g.POST("/connect", Connect(sm))

	v1 := g.Group("/api/v1")
	RegisterPing(v1)

	v1.Use(Auth(sm))

	RegisterTls(v1)
}
