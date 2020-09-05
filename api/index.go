package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Index api index
func Index(c echo.Context) error {
	resp := &RespJSON{
		Code:    http.StatusOK,
		Message: "OK",
		Data:    map[string]string{"Server": "Hyuga"},
		Success: true,
	}
	return c.JSON(http.StatusOK, resp)
}
