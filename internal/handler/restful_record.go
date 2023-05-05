package handler

import (
	"strings"

	"github.com/ac0d3r/hyuga/internal/oob"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (w *restfulHandler) all(c *gin.Context) {
	type Request struct {
		Token  string `json:"token" form:"token" binding:"required"`
		Type   string `json:"type" form:"type"`
		Filter string `json:"filter" form:"filter"`
	}

	param := new(Request)
	if BindParam(c, param) {
		return
	}
	switch param.Type {
	case "dns", "http", "ldap", "rmi", "":
	default:
		ReturnError(c, errParams, nil)
		return
	}

	user, err := w.db.GetUserByToken(param.Token, true)
	if err != nil || user == nil {
		ReturnUnauthorized(c, errUnauthorizedAccess)
		return
	}

	// get user all records
	records, err := w.recorder.Get(user.Sid)
	if err != nil {
		logrus.Warnf("[record] get user records err: %s", err.Error())
		return
	}

	data := make([]any, 0)
	for _, r := range records {
		record, ok := r.(oob.Record)
		if !ok {
			continue
		}

		if param.Type != "" {
			if record.Type.String() != param.Type {
				continue
			}
		}
		if param.Filter != "" {
			if !strings.Contains(record.Name, param.Filter) {
				continue
			}
		}
		data = append(data, record)
	}
	ReturnJSON(c, data)
}
