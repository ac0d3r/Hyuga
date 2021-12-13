package oob

import (
	"hyuga/config"
	"hyuga/database"
	"hyuga/handler/rand"
	"log"
	"net"
	"strings"
	"time"

	"github.com/miekg/dns"
)

const (
	LogTTL     = 0
	NsTTL      = 10 * 60
	DefaultTTL = 5 * 60
)

type DnsServer struct {
	server *dns.Server
}

func NewDnsServer(addr string) *DnsServer {
	server := &DnsServer{}
	server.server = &dns.Server{
		Addr:    addr + ":53",
		Net:     "udp",
		Handler: server,
	}
	return server
}

func (d *DnsServer) ListenAndServe() {
	if err := d.server.ListenAndServe(); err != nil {
		log.Printf("Could not serve dns on port 53: %s\n", err)
	}
}

func (d *DnsServer) Shutdown() error {
	return d.server.Shutdown()
}

func (d *DnsServer) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
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
	name := strings.Trim(question.Name, ".")
	identity := getIdentity(name, config.C.Domain.Main)
	value := ""
	if identity != "" && database.UserExist(identity) {
		rhost, _, _ := net.SplitHostPort(w.RemoteAddr().String())
		record := database.DnsRecord{
			Name:       name,
			RemoteAddr: rhost,
			Created:    time.Now().Unix(),
		}
		if err := database.SetUserRecord(identity, record, config.C.RecordExpiration); err != nil {
			log.Printf("SetUserRecord %s %#v error: %s\n", identity, record, err)
		}
	} else {
		value = config.C.Domain.IP
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
		rrHeader.Ttl = LogTTL
		if value == "" {
			value = getDnsValue(identity)
		}
		rrs = append(rrs, &dns.A{Hdr: rrHeader, A: net.ParseIP(value)})
	case dns.TypeNS:
		rrHeader.Ttl = NsTTL
		for i := range config.C.Domain.NS {
			rrs = append(rrs, &dns.NS{Hdr: rrHeader, Ns: config.C.Domain.NS[i]})
		}
	default:
		dns.HandleFailed(w, r)
		return
	}

	m.Answer = append(m.Answer, rrs...)
	if err := w.WriteMsg(m); err != nil {
		log.Printf("Failed to write message error: %s \n", err)
	}
}

func getDnsValue(identity string) string {
	ips, err := database.GetUserDNSRebinding(identity)
	if err != nil {
		log.Printf("GetUserDNSRebinding error: %s \n", err)
		return config.C.Domain.IP
	}

	if len(ips) == 0 {
		return config.C.Domain.IP
	}
	ips = append(ips, config.C.Domain.IP)
	return ips[rand.RandomInt(0, len(ips)-1)]
}
