package oob

import (
	"hyuga/config"
	"hyuga/database"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func HttpLog(c *gin.Context) {
	host := GetRequestHost(c.Request.Host)
	identity := getIdentity(host, config.C.Domain.Main)

	if identity != "" && database.UserExist(identity) {
		req, _ := httputil.DumpRequest(c.Request, true)
		rip, _, _ := net.SplitHostPort(getRealIP(c.Request))
		record := database.HttpRecord{
			URL:        c.Request.URL.String(),
			Method:     c.Request.Method,
			RemoteAddr: rip,
			Created:    time.Now().Unix(),
			Raw:        string(req),
		}
		if err := database.SetUserRecord(identity, record, config.C.RecordExpiration); err != nil {
			log.Printf("SetUserRecord %s %#v error: %s", identity, record, err)
		}

		c.Status(http.StatusOK)
		c.Writer.Write([]byte(http.StatusText(http.StatusOK)))
		return
	}
	c.Status(http.StatusOK)
}

func getRealIP(r *http.Request) string {
	ip := r.Header.Get("X-Real-Ip")
	if ip == "" {
		ip = r.Header.Get("X-Forwarded-For")
	}
	if ip == "" {
		ip = r.RemoteAddr
	}
	return ip
}

func getIdentity(domain, sub string) string {
	i := strings.Index(domain, sub)
	if i == 0 || i == -1 {
		return ""
	}

	pre := strings.Split(strings.Trim(domain[:i], "."), ".")
	return pre[len(pre)-1]
}

func GetRequestHost(s string) string {
	var host string
	if strings.Contains(s, ":") {
		host, _, _ = net.SplitHostPort(s)
	} else {
		host = s
	}
	return host
}
