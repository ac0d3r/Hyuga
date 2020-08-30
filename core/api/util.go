package api

import (
	"net/http"
)

// RespJSON respsone json struct
type RespJSON struct {
	Code    uint16      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Success bool        `json:"success"`
}

// ProcessingError processing error
func ProcessingError(err error) (int, *RespJSON) {
	return http.StatusBadRequest, &RespJSON{
		Code:    http.StatusBadRequest,
		Message: err.Error(),
		Data:    nil,
		Success: false,
	}
}
