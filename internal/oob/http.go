package oob

import (
	"net"
	"net/http"
	"net/http/httputil"
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
	remote := getRealIP(c.Request)
	url := c.Request.URL.String()

	logrus.Infof("[oob][http] request '%s' from '%s'", url, remote)
	if sid != "" {
		req, _ := httputil.DumpRequest(c.Request, true)
		if err := h.recorder.Record(sid, Record{
			Sid:        sid,
			Type:       TypeHTTP,
			Name:       c.Request.URL.String(),
			RemoteAddr: remote,
			CreatedAt:  time.Now().Unix(),
			Detail:     map[string]string{"raw": string(req)},
		}); err != nil {
			logrus.Warnf("[http] set record %s %#v error: %s", sid, url, err)
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
	host, _, _ := net.SplitHostPort(ip)
	return host
}
