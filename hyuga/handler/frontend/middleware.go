package frontend

import (
	"hyuga/config"
	"hyuga/database"
	"hyuga/handler/base"
	"hyuga/oob"
	"strings"

	"github.com/gin-gonic/gin"
)

func MiddlewareUserToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string
		authorization := c.Request.Header.Get("Authorization")
		if authorization == "" {
			token = c.Query("token")
		} else {
			token = strings.TrimPrefix(authorization, "Bearer ")
		}

		if token == "" {
			base.ReturnUnauthorized(c, 200)
			c.Abort()
			return
		}

		user, err := database.GetUserByToken(token)
		if err != nil {
			base.ReturnUnauthorized(c, 200)
			c.Abort()
			return
		}

		c.Set("uid", user.ID)
		c.Set("token", user.Token)
		c.Next()
	}
}

func MiddlewareForwardLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		host := strings.Split(c.Request.Host, ":")[0]
		if host != config.MainDomain {
			oob.HttpLog(c)
			c.Abort()
		}
	}
}
