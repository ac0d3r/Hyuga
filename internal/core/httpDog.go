package core

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"hyuga/conf"
	"hyuga/db"
	"hyuga/internal"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func fullURL(scheme, host string, req *url.URL) string {
	return fmt.Sprintf("%s://%s%s", scheme, host, req.RequestURI())
}

func splicingCookies(cookies []*http.Cookie) string {
	var c string
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

// HTTPDog Echo middleware record request that do not belong to this API
func HTTPDog() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			host := strings.Split(req.Host, ":")[0]
			log.Debug("httpDog Request host: ", host)
			// if conf.AppENV == "development" {
			// 	return next(c)
			// }
			if host == ("api." + conf.Domain.Main) {
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
				"url":        fullURL(c.Scheme(), host, req.URL),
				"method":     req.Method,
				"remoteAddr": internal.CutStrings(requestedRealIP(req), ":"),
				"cookies":    splicingCookies(req.Cookies()),
			}
			identity := parseIdentity(host)
			log.Debug("httpDog identity: ", identity)
			if identity != "" {
				err := db.SetRecord(ctx, db.HTTPType, identity, httpData)
				if err != nil {
					log.Error("httpDog error: ", err)
					return c.JSON(internal.ProcessingError(err))
				}
				log.Debug("httpDog: ", httpData)
				return c.JSON(http.StatusOK, &internal.RespJSON{
					Code:    http.StatusOK,
					Message: "OK",
					Data:    nil,
					Success: true,
				})
			}
			return nil
		}
	}
}
