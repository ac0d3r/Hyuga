package handler

import (
	"fmt"
	"net/http"

	"github.com/ac0d3r/hyuga/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Register interface {
	RegisterHandler(g *gin.Engine)
}

func BindParam(c *gin.Context, obj any) bool {
	if err := c.ShouldBind(obj); err != nil {
		ReturnError(c, errParams, err)
		return true
	}
	return false
}

type Result struct {
	Code errCode     `json:"code"`
	Msg  string      `json:"message,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

// ReturnJSON 返回正确的业务数据
func ReturnJSON(c *gin.Context, data any) {
	defer func() {
		if e := recover(); e != nil {
			ReturnError(c, errInternal, fmt.Errorf("recover error: %v", e))
		}
	}()
	c.JSON(http.StatusOK, Result{
		Data: data,
	})
}

type errCode int

const (
	errInternal errCode = 1000 + iota
	errParams
	errDatabase

	errOauth
	errUnauthorizedAccess
)

var errorMap = map[errCode]string{
	errInternal:           "服务器内部错误",
	errParams:             "请求参数错误",
	errDatabase:           "数据库操作错误",
	errOauth:              "授权错误",
	errUnauthorizedAccess: "鉴权错误",
}

// ReturnError 返回错误
func ReturnError(c *gin.Context, code errCode, err error) {
	msg, has := errorMap[code]
	if !has {
		msg = "unknown error code"
	}

	if err != nil {
		logrus.Warnf("[handler] {%s}->[%s] error: #(%v)", c.Request.RequestURI, logger.FindCaller(), err)
	}

	c.JSON(http.StatusOK, Result{
		Code: code,
		Msg:  msg,
	})
}

func ReturnUnauthorized(c *gin.Context, code errCode, err ...error) {
	msg, has := errorMap[code]
	if !has {
		msg = "unknown error code"
	}
	if len(err) > 0 {
		if err != nil {
			logrus.Warnf("[handler] {%s}->[%s] error: #(%v)", c.Request.RequestURI, logger.FindCaller(), err)
		}
	}

	c.JSON(http.StatusUnauthorized, Result{
		Code: code,
		Msg:  msg,
	})
}

var cookieExpireTime = 3600 * 24 * 7

func setCookie(c *gin.Context, key, value string) {
	c.SetCookie(key, value, cookieExpireTime, "/", "", true, false)
}
