package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {

		accessToken := c.GetHeader("accessToken")
		fmt.Println("accessToken", accessToken)

		// username := c.PostForm("user")
		// password := c.PostForm("password")

		// if username == "foo" && password == "bar" {
		// 	return
		// } else {
		// 	c.AbortWithStatus(http.StatusUnauthorized)
		// }
	}
}
