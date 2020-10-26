package utils

import (
	"regexp"
	"strings"

	"github.com/labstack/gommon/log"
)

// ParseRemoteAddr remove port
func ParseRemoteAddr(addr, step string) string {
	if strings.Contains(addr, step) {
		return strings.Split(addr, step)[0]
	}
	return addr
}

// StringSlice2AnySlice []string -> []insterface{}
func StringSlice2AnySlice(old []string) []interface{} {
	new := make([]interface{}, len(old))
	for i, v := range old {
		new[i] = v
	}
	return new
}

// CheckIP simple filter ip
func CheckIP(IP string) bool {
	match, err := regexp.MatchString(`\d+\.\d+\.\d+\.\d+`, IP)
	if !match || err != nil {
		log.Error("CheckIP match: ", match, err, " ", IP)
		return false
	}
	return true
}
