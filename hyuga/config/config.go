package config

import (
	"io"
	"os"
	"strings"
	"time"

	"github.com/miekg/dns"
	"gopkg.in/yaml.v2"
)

var (
	C = Config{}
)

type Config struct {
	DebugMode            bool          `yaml:"debug"`
	Redis                string        `yaml:"redis"`
	Domain               domainSetting `yaml:"domain"`
	RecordExpirationHour int           `yaml:"record_expiration_hours"`
	RecordExpiration     time.Duration
}

type domainSetting struct {
	Main string   `yaml:"main"`
	NS   []string `yaml:"ns"`
	IP   string   `yaml:"ip"`
}

func SetFromYaml(c string) error {
	f, err := os.Open(c)
	if err != nil {
		return err
	}

	buf, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	if err = yaml.Unmarshal(buf, &C); err != nil {
		return err
	}

	C.Domain.Main = strings.Trim(C.Domain.Main, ".")
	for i := range C.Domain.NS {
		C.Domain.NS[i] = dns.Fqdn(C.Domain.NS[i])
	}

	C.RecordExpiration = time.Duration(time.Duration(C.RecordExpirationHour) * time.Hour)
	return nil
}
