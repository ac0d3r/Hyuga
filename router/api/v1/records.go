package v1

import (
	"fmt"
	"net/http"

	"hyuga/db"
	"hyuga/internal"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

// GetRecords get records
func GetRecords(c echo.Context) error {
	rtype := c.QueryParam("type")
	token := c.QueryParam("token")
	filter := c.QueryParam("filter")

	log.Debug(fmt.Sprintf("api/v1/GetRecords: type=%s token=%s filter=%s", rtype, token, filter))

	records, err := db.GetRecord(ctx, rtype, token, filter)
	if err != nil {
		return c.JSON(internal.ProcessingError(err))
	}
	internal.SortRecords(records)
	code := http.StatusOK
	return c.JSON(code, &internal.RespJSON{
		Code:    code,
		Message: http.StatusText(code),
		Data:    records,
		Success: true,
	})
}
