package v1

import (
	"Hyuga/api"
	"Hyuga/database"
	"Hyuga/utils"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

// GetRecords get records
func GetRecords(c echo.Context) error {
	rtype := c.QueryParam("type")
	token := c.QueryParam("token")
	filter := c.QueryParam("filter")

	log.Debug(fmt.Sprintf("api/v1/GetRecords: type=%s token=%s filter=%s", rtype, token, filter))

	records, err := database.Recorder.GetRecords(rtype, token, filter)
	if err != nil {
		code, resp := api.ProcessingError(err)
		return c.JSON(code, resp)
	}
	utils.SortRecords(records)
	code := http.StatusOK
	return c.JSON(code, &api.RespJSON{
		Code:    code,
		Message: http.StatusText(code),
		Data:    records,
		Success: true,
	})
}
