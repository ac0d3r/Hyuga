package internal

import (
	"regexp"
	"strings"
)

// CutStrings 移除 `cut` 字符之后的内容
func CutStrings(addr, cut string) string {
	if strings.Contains(addr, cut) {
		return strings.Split(addr, cut)[0]
	}
	return addr
}

// StringSlice2AnySlice []string -> []insterface{}
func StringSlice2AnySlice(old []string) []interface{} {
	news := make([]interface{}, len(old))
	for i, v := range old {
		news[i] = v
	}
	return news
}

// CheckIP simple filter ip
func CheckIP(IP string) (bool, error) {
	match, err := regexp.MatchString(`\d+\.\d+\.\d+\.\d+`, IP)
	if !match || err != nil {
		return false, err
	}
	return true, nil
}
