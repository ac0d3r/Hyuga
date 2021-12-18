package base

import (
	"fmt"
	"hyuga/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Result struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

func ReturnJSON(c *gin.Context, data interface{}) {
	defer func() {
		if r := recover(); r != nil {
			c.JSON(http.StatusOK, gin.H{"code": 1000, "msg": "internal system error"})
		}
	}()
	c.JSON(http.StatusOK, Result{
		Data: data,
	})
}

func ReturnError(c *gin.Context, code int, err ...error) {
	msg, has := errorMap[code]
	if !has {
		msg = "unknown error code"
	}
	if len(err) > 0 && config.DebugMode {
		msg = msg + fmt.Sprintf(" [Error %v]", err[0])
	}

	c.JSON(http.StatusOK, Result{
		Code: code,
		Msg:  msg,
	})
}

func ReturnUnauthorized(c *gin.Context, code int, err ...error) {
	msg, has := errorMap[code]
	if !has {
		msg = "unknown error code"
	}
	if len(err) > 0 && config.DebugMode {
		msg = msg + fmt.Sprintf(" [Error %v]", err[0])
	}

	c.JSON(http.StatusUnauthorized, Result{
		Code: code,
		Msg:  msg,
	})
}

var errorMap = map[int]string{
	100: "internal system error",
	101: "request parameters error",
	102: "database error",
	200: "unauthorized access",
}

func GetUserID(c *gin.Context) string {
	return c.GetString("uid")
}
