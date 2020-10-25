package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Index api index
func Index(c echo.Context) error {
	code := http.StatusOK
	resp := &RespJSON{
		Code:    code,
		Message: http.StatusText(code),
		Data:    map[string]string{"Server": "Hyuga"},
		Success: true,
	}
	return c.JSON(code, resp)
}
