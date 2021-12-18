package config

import (
	"io"
	"net"
	"os"
	"strings"
	"time"

	"github.com/miekg/dns"
	"gopkg.in/yaml.v2"
)

var (
	DebugMode        bool
	RedisDsn         string
	RecordExpiration time.Duration
	MainDomain       string
	DefaultIP        net.IP
	NSDomain         []string
)

type Config struct {
	DebugMode bool   `yaml:"debug"`
	Redis     string `yaml:"redis"`
	Domain    struct {
		Main string   `yaml:"main"`
		NS   []string `yaml:"ns"`
		IP   string   `yaml:"ip"`
	} `yaml:"domain"`
	RecordExpirationHour int `yaml:"record_expiration_hours"`
}

func SetFromYaml(c string) error {
	var conf Config

	f, err := os.Open(c)
	if err != nil {
		return err
	}

	buf, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	if err = yaml.Unmarshal(buf, &conf); err != nil {
		return err
	}

	DebugMode = conf.DebugMode
	RedisDsn = conf.Redis
	RecordExpiration = time.Duration(time.Duration(conf.RecordExpirationHour) * time.Hour)

	MainDomain = strings.Trim(conf.Domain.Main, ".")
	NSDomain = make([]string, len(conf.Domain.NS))
	for i := range conf.Domain.NS {
		NSDomain[i] = dns.Fqdn(conf.Domain.NS[i])
	}
	DefaultIP = net.ParseIP(conf.Domain.IP)

	return nil
}
