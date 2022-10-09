package api

import (
	"fmt"
	"net/http"
	"zkpass-node/pkg/session"

	"github.com/gin-gonic/gin"
)

func RegisterTls(r *gin.RouterGroup) {
	r.GET("/tls", mac)
}

func mac(c *gin.Context) {
	s := c.MustGet("session").(*session.Session)
	fmt.Println("mac" + s.Id)
	c.JSON(http.StatusOK, gin.H{
		"mes": "mac changed",
	})
}
