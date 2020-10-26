package core

import (
	"Hyuga/conf"
	"Hyuga/database"
	"Hyuga/utils"
	"fmt"
	"regexp"
	"strings"

	"github.com/labstack/gommon/log"
	"github.com/miekg/dns"
)

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
	IP, err = database.Recorder.GetUserDNSRebinding(identity)
	return
}

func giveAnswer(identity, qName string, qType uint16) (answers map[string][]string) {
	// +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
	// 	*.hyuga.io.		IN 	NS		ns1.buzz.io.
	// 	*.hyuga.io.		IN 	NS		ns2.buzz.io.
	// 	*.hyuga.io.		IN 	A		1.1.1.1
	// 	hyuga.io. 		IN 	A   	1.1.1.1
	// dnsRebinding
	// 	*.r.*.hyuga.io.	IN 	A		`rebinding hosts`
	// +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
	answers = make(map[string][]string)
	if !strings.HasSuffix(qName, fmt.Sprintf("%s.", conf.Domain)) {
		return
	}
	ips := []string{conf.ServerIP}
	// handler dns-rebinding
	IP := getDNSRebinding(identity, qName)
	if IP != "" {
		ips[0] = IP
	}

	switch qType {
	case dns.TypeA:
		answers["A"] = ips
	case dns.TypeNS:
		answers["NS"] = []string{conf.NS1Domain, conf.NS2Domain}
	case dns.TypeANY:
		answers["A"] = ips
		answers["NS"] = []string{conf.NS1Domain, conf.NS2Domain}
	}
	return
}

func parseQuery(remoteAddr string, m *dns.Msg) {
	for _, q := range m.Question {
		// record dns query
		identity := parseIdentity(q.Name)
		if identity != "" {
			dnsData := map[string]interface{}{
				"name":       strings.TrimRight(q.Name, "."),
				"remoteAddr": remoteAddr}
			err := database.Recorder.Record("dns", identity, dnsData)
			if err != nil {
				log.Error("dnsDog: ", err)
			}
		}

		// make answers for dns query
		answers := giveAnswer(identity, q.Name, q.Qtype)
		for qtype := range answers {
			log.Debug(fmt.Sprintf("dnsDog: Query for %s %s", qtype, q.Name))
			for _, record := range answers[qtype] {
				rr, err := dns.NewRR(fmt.Sprintf("%s %s %s", q.Name, qtype, record))
				if err == nil {
					m.Answer = append(m.Answer, rr)
				}
			}
		}
	}
}

func handleDNSRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	switch r.Opcode {
	case dns.OpcodeQuery:
		parseQuery(utils.ParseRemoteAddr(w.RemoteAddr().String(), ":"), m)
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
