package core

import (
	"context"
	"fmt"
	"net"
	"regexp"
	"strings"

	"hyuga/conf"
	"hyuga/db"
	"hyuga/internal"
	"hyuga/internal/random"

	"github.com/labstack/gommon/log"
	"github.com/miekg/dns"
)

const (
	LogTTL     = 0
	NsTTL      = 10 * 60
	DefaultTTL = 5 * 60
	XIPTTL     = 24 * 60 * 60
)

var (
	ctx                context.Context = context.Background()
	ignoreIdentityList                 = []string{"api", "admin"}
)

func isIgnoreIdentity(identity string) bool {
	for _, v := range ignoreIdentityList {
		if identity == v {
			return true
		}
	}
	return false
}

func parseIdentity(domain string) string {
	reg := regexp.MustCompile(fmt.Sprintf(`\.?([^\.]+)\.%s\.?`, conf.Domain.Main))
	subs := reg.FindStringSubmatch(domain)
	if len(subs) >= 2 {
		return subs[1]
	}
	return ""
}

func getDNSRebinding(identity, qName string) (IP string) {
	if identity == "" {
		return
	}
	match, err := regexp.MatchString(fmt.Sprintf(`\.?.*r\.%s\.%s\.?`, identity, conf.Domain.Main), qName)
	if !match || err != nil {
		log.Debug("getDNSRebinding regexp match: ", match, err)
		return
	}
	ips, err := db.GetUserDNSRebinding(ctx, identity, true)
	if err != nil {
		log.Error("getDNSRebinding", err)
	}
	if len(ips) >= 1 {
		IP = ips[0]
	}
	return
}

func giveAnswer(identity, qName string, qType uint16, ttl uint32) dns.RR {
	// +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
	// 	*.hyuga.io.		IN 	NS		ns1.buzz.io.
	// 	*.hyuga.io.		IN 	NS		ns2.buzz.io.
	// 	*.hyuga.io.		IN 	A		1.1.1.1
	// 	hyuga.io. 		IN 	A   	1.1.1.1
	// dnsRebinding
	// 	*.r.*.hyuga.io.	IN 	A		`rebinding hosts`
	// +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
	if !strings.HasSuffix(qName, fmt.Sprintf("%s.", conf.Domain.Main)) {
		return nil
	}
	respip := conf.Domain.IP
	// handler dns-rebinding
	if ttl == LogTTL {
		rebindingIP := getDNSRebinding(identity, qName)
		if rebindingIP != "" {
			respip = rebindingIP
		}
	}
	rrHeader := dns.RR_Header{
		Name:   qName,
		Rrtype: qType,
		Class:  dns.ClassINET,
		Ttl:    ttl,
	}
	switch qType {
	case dns.TypeA:
		return &dns.A{Hdr: rrHeader, A: net.ParseIP(respip)}
	case dns.TypeNS:
		return &dns.NS{Hdr: rrHeader, Ns: conf.Domain.NS[random.Int(0, len(conf.Domain.NS)-1)]}
	default:
		return nil
	}
}

func parseQuery(w dns.ResponseWriter, req *dns.Msg, m *dns.Msg) {
	ttl := DefaultTTL
	for _, q := range m.Question {
		if q.Qclass != dns.ClassINET {
			dns.HandleFailed(w, req)
			return
		}
		identity := parseIdentity(q.Name)
		// record dns query
		if identity != "" && !isIgnoreIdentity(identity) {
			dnsData := map[string]interface{}{
				"name":       strings.TrimRight(q.Name, "."),
				"remoteAddr": internal.CutStrings(w.RemoteAddr().String(), ":")}
			err := db.SetRecord(ctx, db.DNSType, identity, dnsData)
			if err != nil {
				log.Error("dnsDog: ", err)
			} else {
				ttl = LogTTL
			}
		}

		// make answers for dns query
		log.Debug(fmt.Sprintf("dnsDog: Query for %s %s",
			dns.Type(q.Qtype).String(),
			q.Name))

		switch q.Qtype {
		case dns.TypeANY:
			fallthrough
		case dns.TypeA:
			a := giveAnswer(identity, q.Name, dns.TypeA, uint32(ttl))
			if a != nil {
				m.Answer = append(m.Answer, a)
			}
		case dns.TypeAAAA:
			// not ipv6 now
			dns.HandleFailed(w, req)
		case dns.TypeNS:
			ttl = NsTTL
			a := giveAnswer(identity, q.Name, dns.TypeNS, uint32(ttl))
			if a != nil {
				m.Answer = append(m.Answer, a)
			}
		default:
			dns.HandleFailed(w, req)
		}
	}
}

func handleDNSRequest(w dns.ResponseWriter, req *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(req)
	m.Compress = false

	switch req.Opcode {
	case dns.OpcodeQuery:
		parseQuery(w, req, m)
	}
	if err := w.WriteMsg(m); err != nil {
		log.Error(fmt.Sprintf("Failed to write message:%s", err))
	}
}

// NewDNSDog dnsDog serve
func NewDNSDog(addr string) (*dns.Server, error) {
	dns.HandleFunc(conf.Domain.Main+".", handleDNSRequest)
	return &dns.Server{Addr: addr, Net: "udp"}, nil
}
