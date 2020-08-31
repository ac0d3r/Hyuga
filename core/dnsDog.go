package core

import (
	"Hyuga/core/conf"
	"Hyuga/core/utils"
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

func giveAnswer(qName string, qType uint16) (answers map[string][]string) {
	// +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
	// 	*.hyuga.io.		IN 	NS		ns1.buzz.io.
	// 	*.hyuga.io.		IN 	NS		ns2.buzz.io.
	// 	*.hyuga.io.		IN 	A		1.1.1.1
	// 	hyuga.io. 		IN 	A   	1.1.1.1
	// +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
	answers = make(map[string][]string)
	if !strings.HasSuffix(qName, fmt.Sprintf("%s.", conf.Domain)) {
		return
	}
	multRe := fmt.Sprintf(`.*\.%s\.`, conf.Domain)

	switch qType {
	case dns.TypeA:
		answers["A"] = []string{conf.ServerIP}
	case dns.TypeNS:
		match, _ := regexp.MatchString(multRe, qName)
		if match {
			answers["NS"] = []string{conf.NS1Domain, conf.NS2Domain}
		}
	case dns.TypeANY:
		answers["A"] = []string{conf.ServerIP}
		answers["NS"] = []string{conf.NS1Domain, conf.NS2Domain}
	}
	return
}

func parseQuery(remoteAddr string, m *dns.Msg) {
	records := false
	for _, q := range m.Question {
		// only record dns queries once
		if !records {
			dnsData := map[string]interface{}{"name": q.Name, "remoteAddr": remoteAddr}
			identity := parseIdentity(q.Name)
			if identity == "" {
				goto ANSWER
			}
			err := utils.Recorder.Record("dns", identity, dnsData)
			if err != nil {
				log.Error("dnsDog: ", err)
			} else {
				records = true
			}
		}

	ANSWER:
		answers := giveAnswer(q.Name, q.Qtype)
		for qtype := range answers {
			log.Debug(fmt.Sprintf("Query for %s %s", qtype, q.Name))
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
		parseQuery(w.RemoteAddr().String(), m)
	}

	err := w.WriteMsg(m)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to write message:%s", err))
	}
}

// DNSDogServe dnsDog serve
func DNSDogServe() {
	// attach request handler func
	dns.HandleFunc(fmt.Sprintf("%s.", conf.Domain), handleDNSRequest)

	// start server
	addr := ":53"
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
