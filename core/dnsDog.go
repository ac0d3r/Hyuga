package core

import (
	"Hyuga/conf"
	"Hyuga/database"
	"Hyuga/utils"
	"fmt"
	"net"
	"regexp"
	"strings"

	"github.com/labstack/gommon/log"
	"github.com/miekg/dns"
)

const (
	LogTTL     = 0
	NsTTL      = 10 * 60
	DefaultTTL = 5 * 60
	XIPTTL     = 24 * 60 * 60
)

var ignoreIdentityList = []string{"api", "admin"}
var nss = []string{conf.NS1Domain, conf.NS2Domain}

func isIgnoreIdentity(identity string) bool {
	for _, v := range ignoreIdentityList {
		if identity == v {
			return true
		}
	}
	return false
}

func parseIdentity(domain string) string {
	reg := regexp.MustCompile(fmt.Sprintf(`\.?([^\.]+)\.%s\.?`, conf.Domain))
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
	match, err := regexp.MatchString(fmt.Sprintf(`\.?.*r\.%s\.%s\.?`, identity, conf.Domain), qName)
	if !match || err != nil {
		log.Debug("getDNSRebinding regexp match: ", match, err)
		return
	}
	IPs, err := database.Recorder.GetUserDNSRebinding(identity, true)
	if err != nil {
		log.Error("getDNSRebinding", err)
	}
	if len(IPs) >= 1 {
		IP = IPs[0]
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
	if !strings.HasSuffix(qName, fmt.Sprintf("%s.", conf.Domain)) {
		return nil
	}
	respip := conf.ServerIP
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
		return &dns.NS{Hdr: rrHeader, Ns: nss[utils.RandInt(0, len(nss)-1)]}
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
				"remoteAddr": utils.ParseRemoteAddr(w.RemoteAddr().String(), ":")}
			err := database.Recorder.Record("dns", identity, dnsData)
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
	err := w.WriteMsg(m)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to write message:%s", err))
	}
}

// DNSDogServe dnsDog serve
func DNSDogServe(addr string) {
	// attach request handler func
	dns.HandleFunc(fmt.Sprintf("%s.", conf.Domain), handleDNSRequest)

	server := &dns.Server{Addr: addr, Net: "udp"}
	log.Info(fmt.Sprintf("Starting at %s", addr))
	err := server.ListenAndServe()
	defer func() {
		err := server.Shutdown()
		if err != nil {
			log.Error(err, "Failed to shutdown fakedns server")
		}

	}()
	if err != nil {
		log.Error(err, "Failed to start fakedns server")
	}
}
