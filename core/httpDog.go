package core

import (
	"Hyuga/api"
	"Hyuga/conf"
	"Hyuga/database"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func splicingCookies(cookies []*http.Cookie) string {
	c := ""
	for _, cookie := range cookies {
		c += cookie.String()
	}
	return c
}

func requestedRealIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}

// httpDog Echo middleware record request that do not belong to this API
func httpDog() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			host := strings.Split(req.Host, ":")[0]
			log.Debug("httpDog Request host: ", host)
			// release requests from `api.huyga.io`
			// if conf.AppENV == "development" {
			// 	return next(c)
			// }
			if host == fmt.Sprintf("api.%s", conf.Domain) {
				return next(c)
			}

			// record other requests from `*.huyga.io`
			// +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
			// url
			// method
			// remoteAddr
			// cookies
			// +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
			httpData := map[string]interface{}{
				"url":        req.RequestURI,
				"method":     req.Method,
				"remoteAddr": requestedRealIP(req),
				"cookies":    splicingCookies(req.Cookies()),
			}
			identity := parseIdentity(host)
			log.Debug("httpDog identity: ", identity)
			if identity != "" {
				err := database.Recorder.Record("http", identity, httpData)
				if err != nil {
					log.Error("httpDog: ", err)
					goto NOTFOUND
				}
				log.Debug("httpDog: ", httpData)
				return c.JSON(http.StatusOK, &api.RespJSON{
					Code:    http.StatusOK,
					Message: "OK",
					Data:    nil,
					Success: true,
				})
			}
		NOTFOUND:
			return c.JSON(http.StatusNotFound, &api.RespJSON{
				Code:    http.StatusNotFound,
				Message: http.StatusText(http.StatusNotFound),
				Data:    nil,
				Success: false,
			})
		}
	}
}
