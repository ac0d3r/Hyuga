package v1

import (
	"Hyuga/core/api"
	"Hyuga/core/conf"
	"Hyuga/core/utils"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func genUser() (identity string, token string) {
	length := 4
	times := 0

	for {
		identity = utils.RandomString(length)
		if !utils.Recorder.IdentityExist(identity) {
			token = utils.RandomString(utils.RandInt(32, 64))
			if !utils.Recorder.UserExist(identity, token) {
				break
			}
		}
		if times < 3 {
			times++
		} else {
			times = 0
			length++
		}
	}
	return
}

// CreateUser create user
func CreateUser(c echo.Context) error {
	identity, token := genUser()
	log.Debug("api/v1/CreateUser ", identity, " / ", token)
	err := utils.Recorder.CreateUser(identity, token)
	if err != nil {
		code, resp := api.ProcessingError(err)
		return c.JSON(code, resp)
	}
	return c.JSON(http.StatusOK, &api.RespJSON{
		Code:    http.StatusOK,
		Message: "OK",
		Data:    map[string]string{"identity": fmt.Sprintf("%s.%s", identity, conf.Domain), "token": token},
		Success: true,
	})
}
