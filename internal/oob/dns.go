package oob

import (
	"context"
	"net"
	"strings"

	"hyuga/internal/config"
	"hyuga/internal/db"

	"github.com/miekg/dns"
	"github.com/sirupsen/logrus"
)

const (
	ZeroTTL    = 0
	DefaultTTL = 60
	NsTTL      = 24 * 60 * 60
)

type Dns struct {
	cnf *config.OOB
	db  *db.Client

	*dns.Server
}

func NewDns(cnf *config.OOB,
	db *db.Client) *Dns {

	d := &Dns{cnf: cnf, db: db}
	d.Server = &dns.Server{
		Addr:    ":53",
		Net:     "udp",
		Handler: d,
	}
	return d
}

func (d *Dns) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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
	sid := parseSid(dnsName, d.cnf.Dns.Domain)

	if sid != "" {
		remoteAddr := strings.Split(w.RemoteAddr().String(), ":")[0]
		if _, err := d.db.CreateDNSRecord(ctx, sid, dnsName, remoteAddr); err != nil {
			logrus.Warnf("[dns] set record '%s' '%#v' error: %s\n", sid, dnsName, err)
		}
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
	case dns.TypeA:
		rrs = append(rrs, &dns.A{Hdr: rrHeader, A: net.IP(d.cnf.Dns.IP)})
	case dns.TypeNS:
		rrHeader.Ttl = NsTTL
		for _, ns := range d.cnf.Dns.NS {
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
