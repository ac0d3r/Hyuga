package base

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

var Debug bool

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
	if len(err) > 0 && Debug {
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
	if len(err) > 0 && Debug {
		msg = msg + fmt.Sprintf(" [Error %v]", err[0])
	}

	c.JSON(http.StatusUnauthorized, Result{
		Code: code,
		Msg:  msg,
	})
}

func BindValidate(c *gin.Context, param interface{}) bool {
	if err := c.ShouldBind(param); err != nil && err != io.EOF {
		ReturnError(c, 1001, err)
		return true
	}
	return false
}

var errorMap = map[int]string{
	1000: "internal system error",
	1001: "request parameters error",
	1002: "database error",
	2000: "unauthorized access",
}

func GetUserID(c *gin.Context) int {
	return c.GetInt("uid")
}
