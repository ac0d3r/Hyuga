package oob

import (
	"hyuga/config"
	"hyuga/database"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func HttpLog(c *gin.Context) {
	host := strings.Split(c.Request.Host, ":")[0]
	identity := getIdentity(host, config.MainDomain)

	if identity != "" && database.UserExist(identity) {
		req, _ := httputil.DumpRequest(c.Request, true)
		record := database.HttpRecord{
			URL:        c.Request.URL.String(),
			Method:     c.Request.Method,
			RemoteAddr: getRealIP(c.Request),
			Created:    time.Now().Unix(),
			Raw:        string(req),
		}
		if err := database.SetUserRecord(identity, record, config.RecordExpiration); err != nil {
			log.Printf("[http] set record %s %#v error: %s", identity, record, err)
		}
	}
	c.Status(http.StatusOK)
	c.Writer.Write([]byte(http.StatusText(http.StatusOK)))
}

func getRealIP(r *http.Request) string {
	ip := r.Header.Get("X-Real-IP")
	if ip == "" {
		ip = r.Header.Get("X-Forwarded-For")
	}
	if ip == "" {
		ip = r.RemoteAddr
	}
	return strings.Split(ip, ":")[0]
}

func getIdentity(domain, sub string) string {
	i := strings.Index(domain, sub)
	if i <= 0 {
		return ""
	}

	pre := strings.Split(strings.Trim(domain[:i], "."), ".")
	return pre[len(pre)-1]
}
