package v1

import (
	"Hyuga/api"
	"Hyuga/conf"
	"Hyuga/database"
	"Hyuga/utils"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func genUser() (identity string, token string) {
	// remove UpperCharset
	c := utils.Charset + utils.Number
	length := 4
	times := 0

	for {
		identity = utils.RandomStringWithCharset(length, c)
		if !database.Recorder.IdentityExist(identity) {
			token = utils.RandomString(utils.RandInt(32, 64))
			if !database.Recorder.UserExist(identity, token) {
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
	err := database.Recorder.CreateUser(identity, token)
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
