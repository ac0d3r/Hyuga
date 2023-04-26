package oob

import (
	"net"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"github.com/ac0d3r/hyuga/internal/config"
	"github.com/ac0d3r/hyuga/internal/record"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type HTTP struct {
	cnf      *config.DNS
	recorder *record.Recorder
}

func NewHTTP(cnf *config.DNS, recorder *record.Recorder) *HTTP {
	return &HTTP{cnf: cnf, recorder: recorder}
}

func (h *HTTP) Record(c *gin.Context) {
	host, _, _ := net.SplitHostPort(c.Request.Host)
	sid := parseSid(host, h.cnf.Main)

	logrus.Infof("[oob][http] query '%s' from '%s'", host, getRealIP(c.Request))

	if sid != "" {
		req, _ := httputil.DumpRequest(c.Request, true)
		if err := h.recorder.Record(sid, OOBRecord{
			Type:       OOBHTTP,
			Name:       c.Request.URL.String(),
			RemoteAddr: getRealIP(c.Request),
			CreatedAt:  time.Now().Unix(),
			Detail:     map[string]string{"raw": string(req)},
		}); err != nil {
			logrus.Warnf("[http] set record %s %#v error: %s", sid, c.Request.URL.String(), err)
		}
	}

	c.String(http.StatusOK, http.StatusText(http.StatusOK))
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
