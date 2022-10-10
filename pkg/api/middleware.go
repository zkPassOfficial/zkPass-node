package api

import (
	"net/http"
	"zkpass-node/pkg/session"

	"github.com/gin-gonic/gin"
)

func Auth(sm *session.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {

		sid := c.GetHeader("sid")

		if sm.Has(sid) {
			session := sm.GetSession(sid)
			c.Set("session", session)
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}
