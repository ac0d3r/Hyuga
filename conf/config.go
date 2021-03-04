package conf

import (
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/labstack/gommon/log"
	yaml "gopkg.in/yaml.v2"
)

type app struct {
	Env                 string `yaml:"env"`
	RecordExpirationDay int    `yaml:"recordExpirationDays"`
}
type domain struct {
	Main string   `yaml:"main"`
	NS   []string `yaml:"ns"`
	IP   string   `yaml:"ip"`
}
type config struct {
	App    app    `yaml:"app"`
	Redis  string `yaml:"redis"`
	Domain domain `yaml:"domain"`
}

var (
	conf config = config{}
)

var (
	// AppEnv export conf.App.Env
	AppEnv string
	// RedisAddr export conf.Redis
	RedisAddr string
	// Domain export conf.Domain
	Domain domain = domain{}
	// RecordExpiration 记录的过期时间
	RecordExpiration time.Duration
	// LogMap logmap
	LogMap = map[string]log.Lvl{"development": log.DEBUG, "production": log.INFO}
)

func setLog(env string) {
	// set level
	if level, ok := LogMap[env]; ok {
		log.SetLevel(level)
	} else {
		log.SetLevel(log.DEBUG)
	}
	// set header
	log.SetHeader(`[${level}] ${time_rfc3339_nano} ${prefix} ` +
		`file:'${short_file}'@${line}:`)
}

// setNsDomain 设置 ns 保证以 `.` 结尾
func setNsDomain(ns string) string {
	length := len(ns)
	if length == 0 {
		return ns
	}
	if strings.HasSuffix(ns, ".") {
		return ns
	}
	return ns + "."
}

// SetFromFile 从文件设置 config
func SetFromFile(c string) error {
	var (
		f   *os.File
		buf []byte
		err error
	)
	if f, err = os.Open(c); err != nil {
		return err
	}
	if buf, err = ioutil.ReadAll(f); err != nil {
		return err
	}

	if err = yaml.Unmarshal(buf, &conf); err != nil {
		return err
	}
	setLog(conf.App.Env)
	AppEnv = conf.App.Env
	RedisAddr = conf.Redis
	Domain = conf.Domain
	RecordExpiration = time.Duration(conf.App.RecordExpirationDay*24*60*60) * time.Second
	return nil
}
