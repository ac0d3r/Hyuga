package oob

import (
	"context"
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
	g.Go(func() error {
		udpServ := &dns.Server{
			Addr:    ":53",
			Net:     "udp",
			Handler: dns_,
		}
		logrus.Infof("[oob] dns is listen at '[%s]%s'", udpServ.Net, udpServ.Addr)
		return udpServ.ListenAndServe()
	})
	g.Go(func() error {
		tcpServ := &dns.Server{
			Addr:    ":53",
			Net:     "tcp",
			Handler: dns_,
		}
		logrus.Infof("[oob] dns is listen at '[%s]%s'", tcpServ.Net, tcpServ.Addr)
		return tcpServ.ListenAndServe()
	})

	jndi := NewJNDI(db, &cnf.JNDI, r)
	g.Go(func() error {
		logrus.Infof("[oob] jndi listen at '%s'", cnf.JNDI.Address)
		return jndi.Run(ctx)
	})
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
