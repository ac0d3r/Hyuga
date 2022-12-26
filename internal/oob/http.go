package oob

import (
	"net/http"
	"net/http/httputil"
	"strings"

	"hyuga/internal/config"
	"hyuga/internal/db"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type HTTPLog struct {
	cnf *config.OOB
	db  *db.Client
}

func NewHTTPLog(cnf *config.OOB, db *db.Client) *HTTPLog {
	return &HTTPLog{
		cnf: cnf,
		db:  db,
	}
}

func (h *HTTPLog) Record(c *gin.Context) {
	host := strings.Split(c.Request.Host, ":")[0]
	sid := parseSid(host, h.cnf.Dns.Domain)

	if sid != "" {
		req, _ := httputil.DumpRequest(c.Request, true)
		if _, err := h.db.CreateHTTPRecord(c.Request.Context(), sid, c.Request.URL.String(), c.Request.Method, getRealIP(c.Request), string(req)); err != nil {
			logrus.Warnf("[http] set record %s %#v error: %s", sid, c.Request, err)
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
