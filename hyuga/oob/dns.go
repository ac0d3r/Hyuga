package oob

import (
	"fmt"
	"hyuga/config"
	"hyuga/database"
	"log"
	"net"
	"strings"
	"time"

	"github.com/miekg/dns"
)

const (
	ZeroTTL    = 0
	DefaultTTL = 60
	NsTTL      = 24 * 60 * 60
)

type DnsServer struct {
	server *dns.Server
}

func NewDnsServer(addr string) *DnsServer {
	server := &DnsServer{}
	server.server = &dns.Server{
		Addr:    addr,
		Net:     "udp",
		Handler: server,
	}
	return server
}

func (d *DnsServer) ListenAndServe() {
	log.Printf("[dns] listen on '%s'\n", d.server.Addr)

	if err := d.server.ListenAndServe(); err != nil {
		log.Printf("[dns] listen fail error: %s\n", err)
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
	identity := getIdentity(name, config.MainDomain)

	var (
		recordtimes    int64
		isDnsRebinding bool
	)
	if identity != "" && database.UserExist(identity) {
		record := database.DnsRecord{
			Name:       name,
			RemoteAddr: strings.Split(w.RemoteAddr().String(), ":")[0],
			Created:    time.Now().Unix(),
		}
		if err := database.SetUserRecord(identity, record, config.RecordExpiration); err != nil {
			log.Printf("[dns] set record '%s' '%#v' error: %s\n", identity, record, err)
		} else {
			if name == fmt.Sprintf("r.%s.%s", identity, config.MainDomain) {
				isDnsRebinding = true
				t, err := database.SetUserDnsRebindingTimes(identity)
				if err != nil {
					log.Printf("[dns] set query times '%s' error: %s\n", identity, err)
				}
				recordtimes = t
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
	case dns.TypeA:
		if isDnsRebinding {
			rrHeader.Ttl = ZeroTTL
		}
		rrs = append(rrs, &dns.A{Hdr: rrHeader, A: getDnsValue(!isDnsRebinding, identity, recordtimes)})
	case dns.TypeNS:
		rrHeader.Ttl = NsTTL
		for i := range config.NSDomain {
			rrs = append(rrs, &dns.NS{Hdr: rrHeader, Ns: config.NSDomain[i]})
		}
	default:
		dns.HandleFailed(w, r)
		return
	}

	m.Answer = append(m.Answer, rrs...)
	if err := w.WriteMsg(m); err != nil {
		log.Printf("[dns] write message fail error: %s \n", err)
	}
}

func getDnsValue(defaults bool, identity string, recordtimes int64) net.IP {
	if defaults {
		return config.DefaultIP
	}

	ips, err := database.GetUserDNSRebinding(identity)
	if err != nil {
		log.Printf("[dns] get DNS-Rebind fail error: %s \n", err)
		return config.DefaultIP
	}

	if len(ips) <= 0 {
		return config.DefaultIP
	}

	ipp := []string{config.DefaultIP.String()}
	ipp = append(ipp, ips...)
	return net.ParseIP(ipp[recordtimes%int64(len(ipp))])
}
