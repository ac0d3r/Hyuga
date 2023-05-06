package oob

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/ac0d3r/hyuga/internal/config"
	"github.com/ac0d3r/hyuga/internal/db"
	"github.com/ac0d3r/hyuga/internal/record"
	"github.com/miekg/dns"
	"github.com/sirupsen/logrus"
)

const (
	ZeroTTL    = 0
	DefaultTTL = 60
	NsTTL      = 24 * 60 * 60
)

type Dns struct {
	db       *db.DB
	cnf      *config.DNS
	recorder *record.Recorder
}

func NewDns(db *db.DB, cnf *config.DNS, recorder *record.Recorder) *Dns {
	return &Dns{db: db, cnf: cnf, recorder: recorder}
}

func (d *Dns) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
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
	logrus.Infof("[oob][dns] query '%s' from '%s'", question.Name, w.RemoteAddr().String())

	var user *db.User
	dnsName := strings.Trim(question.Name, ".")
	sid := parseSid(dnsName, d.cnf.Main)
	if sid != "" {
		u, err := d.db.GetUserBySid(sid)
		if err == nil && u != nil {
			user = u
			ip, _ := SplitHostPort(w.RemoteAddr().String())
			if err := d.recorder.Record(sid, Record{
				Sid:        sid,
				Type:       TypeDNS,
				Name:       dnsName,
				RemoteAddr: ip,
				CreatedAt:  time.Now().Unix(),
			}); err != nil {
				logrus.Warnf("[dns] record sid '%s' error: %s", sid, err)
			}
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
	case dns.TypeA, dns.TypeAAAA:
		if user != nil && d.isDnsRebindQuery(sid, dnsName) {
			dnsa := &dns.A{Hdr: rrHeader, A: net.ParseIP(d.cnf.IP)}
			if len(user.DnsRebind.DNS) > 0 {
				rrHeader.Ttl = ZeroTTL
				dnss := []string{d.cnf.IP}
				dnss = append(dnss, user.DnsRebind.DNS...)
				dnsa.A = net.ParseIP(dnss[user.DnsRebind.Times%int64(len(dnss))])
				// update rebinding times
				user.DnsRebind.Times++
				if err := d.db.UpdateUser(user); err != nil {
					logrus.Warnf("[dns] update user '%s' error: %s", user.Sid, err)
				}
			}
			rrs = append(rrs, dnsa)
		} else {
			rrs = append(rrs, &dns.A{Hdr: rrHeader, A: net.ParseIP(d.cnf.IP)})
		}
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

func (d *Dns) isDnsRebindQuery(sid, domain string) bool {
	return domain == fmt.Sprintf("r.%s.%s", sid, d.cnf.Main)
}
