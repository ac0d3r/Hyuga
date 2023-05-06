package oob

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/ac0d3r/hyuga/internal/config"
	"github.com/ac0d3r/hyuga/internal/db"
	"github.com/ac0d3r/hyuga/internal/record"
	"github.com/miekg/dns"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

func Run(ctx context.Context,
	g *errgroup.Group,
	db *db.DB,
	cnf *config.OOB,
	r *record.Recorder) {

	dns_ := NewDns(db, &cnf.DNS, r)
	udpServ := &dns.Server{
		Addr:    ":53",
		Net:     "udp",
		Handler: dns_,
	}
	tcpServ := &dns.Server{
		Addr:    ":53",
		Net:     "tcp",
		Handler: dns_,
	}

	g.Go(func() error {
		<-ctx.Done()
		logrus.Info("[dns] shutdown")
		var err error
		if e := udpServ.Shutdown(); e != nil {
			err = fmt.Errorf("dns-udp shutdown err:%s", e.Error())
		}
		if e := tcpServ.Shutdown(); e != nil {
			err = fmt.Errorf("dns-tcp shutdown err:%s", e.Error())
		}
		return err
	})
	g.Go(func() error {
		logrus.Infof("[oob] dns is listen at '[%s]%s'", udpServ.Net, udpServ.Addr)
		return udpServ.ListenAndServe()
	})
	g.Go(func() error {
		logrus.Infof("[oob] dns is listen at '[%s]%s'", tcpServ.Net, tcpServ.Addr)
		return tcpServ.ListenAndServe()
	})

	jndi := NewJNDI(db, &cnf.JNDI, r)
	logrus.Infof("[oob] jndi listen at '%s'", cnf.JNDI.Address)
	jndi.Run(ctx, g)
}

const (
	TypeDNS RecordType = iota
	TypeHTTP
	TypeLDAP
	TypeRMI
)

type RecordType uint

func (o RecordType) String() string {
	switch o {
	case TypeDNS:
		return "dns"
	case TypeHTTP:
		return "http"
	case TypeLDAP:
		return "ldap"
	case TypeRMI:
		return "rmi"
	default:
		return "unknown"
	}
}

type Record struct {
	Sid        string            `json:"sid"`
	Type       RecordType        `json:"type"`
	Name       string            `json:"name"`
	RemoteAddr string            `json:"remote_addr"`
	CreatedAt  int64             `json:"created_at"`
	Detail     map[string]string `json:"detail"`
}

func parseSid(domain, mainDomain string) string {
	i := strings.Index(domain, mainDomain)
	if i <= 0 {
		return ""
	}

	pre := strings.Split(strings.Trim(domain[:i], "."), ".")
	return pre[len(pre)-1]
}

func SplitHostPort(host string) (h, p string) {
	i := strings.Index(host, ":")
	if i < 0 {
		return host, ""
	}
	h, p, _ = net.SplitHostPort(host)
	return
}
