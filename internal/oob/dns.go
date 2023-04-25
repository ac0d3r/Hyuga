package oob

import (
	"context"
	"net"
	"strings"

	"github.com/ac0d3r/hyuga/internal/config"
	"github.com/miekg/dns"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

const (
	ZeroTTL    = 0
	DefaultTTL = 60
	NsTTL      = 24 * 60 * 60
)

type Dns struct {
	cnf *config.DNS
}

func NewDns(cnf *config.DNS) *Dns {
	return &Dns{cnf: cnf}
}

func (d *Dns) Run(ctx context.Context, g *errgroup.Group) {
	g.Go(func() error {
		udpServ := &dns.Server{
			Addr:    ":53",
			Net:     "udp",
			Handler: d,
		}
		return udpServ.ListenAndServe()
	})
	g.Go(func() error {
		tcpServ := &dns.Server{
			Addr:    ":53",
			Net:     "tcp",
			Handler: d,
		}
		return tcpServ.ListenAndServe()
	})
}

func (d *Dns) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false
	m.Authoritative = true

	if len(r.Question) == 0 {
		return
	}

	if r.Opcode != dns.OpcodeQuery {
		return
	}

	question := r.Question[0]
	dnsName := strings.Trim(question.Name, ".")
	sid := parseSid(dnsName, d.cnf.Main)

	if sid != "" {
		// TODO
	}

	rrs := make([]dns.RR, 0)
	rrHeader := dns.RR_Header{
		Name:   question.Name,
		Rrtype: question.Qtype,
		Class:  dns.ClassINET,
		Ttl:    DefaultTTL,
	}
	switch question.Qtype {
	case dns.TypeANY:
		fallthrough
	case dns.TypeA, dns.TypeAAAA:
		rrs = append(rrs, &dns.A{Hdr: rrHeader, A: net.IP(d.cnf.IP)})
	case dns.TypeNS:
		rrHeader.Ttl = NsTTL
		for _, ns := range d.cnf.NS {
			rrs = append(rrs, &dns.NS{Hdr: rrHeader, Ns: ns})
		}
	default:
		dns.HandleFailed(w, r)
		return
	}

	m.Answer = append(m.Answer, rrs...)
	if err := w.WriteMsg(m); err != nil {
		logrus.Warnf("[dns] write message fail error: %s \n", err)
	}
}
