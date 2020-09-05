package conf

import (
	"os"
	"strconv"
	"time"

	"strings"

	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
)

// checkNsDomain check ns domain last character "."
func checkNsDomain(ns string) string {
	length := len(ns)
	if length == 0 {
		return ns
	}
	if strings.HasSuffix(ns, ".") {
		return ns
	}
	return ns + "."
}

var (
	// AppENV env "production" or "development"
	AppENV string = "development"
	// RedisHost redis host
	RedisHost string = "localhost"
	// RedisPort redis port
	RedisPort int = 6379
	// RecordExpiration record expiration (7days)
	RecordExpiration time.Duration = time.Second * 60 * 60 * 24 * 7

	// Domain 记录域名&主域名
	Domain string = "hyuga.io"
	// NS1Domain NS1 服务器域名
	NS1Domain string = "ns1.buzz.io"
	// NS2Domain NS2 服务器域名
	NS2Domain string = "ns1.buzz.io"
	// ServerIP 公网IP(A 记录指向IP)
	ServerIP string = "127.0.0.1"
	// LogLevel mapping `AppENV` => log.level
	LogLevel = map[string]log.Lvl{"development": log.DEBUG, "production": log.INFO}
)

func init() {
	var err error
	loadEnv()
	AppENV = os.Getenv("APP_ENV")
	RedisHost = os.Getenv("REDIS_SERVER")
	RedisPort, err = strconv.Atoi(os.Getenv("REDIS_PORT"))
	if err != nil {
		RedisPort = 6379
	}
	Domain = os.Getenv("DOMAIN")
	NS1Domain = checkNsDomain(os.Getenv("NS1_DOMAIN"))
	NS2Domain = checkNsDomain(os.Getenv("NS2_DOMAIN"))
	ServerIP = os.Getenv("SERVER_IP")
	setLog()
}

// loadEnv load environment from '.env' file
func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Error("Error loading .env file")
	}
}

func setLog() {
	// set level
	if level, ok := LogLevel[AppENV]; ok {
		log.SetLevel(level)
	} else {
		// the others are DEBUG(1)
		log.SetLevel(log.DEBUG)
	}
	// set header
	header := `[${level}] ${time_rfc3339_nano} ${prefix} ` +
		`file:'${short_file}'@${line}:`
	log.SetHeader(header)
}
