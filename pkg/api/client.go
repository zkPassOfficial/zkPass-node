package api

import (
	"fmt"
	"net/http"
	"zkpass-node/pkg/session"

	"github.com/gin-gonic/gin"
)

//curl -H "Content-Type: application/json" -X POST -d '{"name": "gin"}'
func Connect(sm *session.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.PostForm("name")
		isValid := true
		if !isValid {
			fmt.Println("Error: invalid " + name)
			// TODO: for Session with new rotated of key, needs to recreate it
			c.JSON(http.StatusOK, gin.H{
				"mes": "name:" + name,
			})
		} else {
			fmt.Println("Login Successful: name " + name)

			s := new(session.Session)
			// c.destroyChan = cm.destroyChan
			sm.Add(s)

			c.Set("sid:"+s.Id, s)

			c.JSON(http.StatusOK, gin.H{
				"sid": s.Id,
			})
		}
	}
}
