package utils

import "strings"

// ParseRemoteAddr remove port
func ParseRemoteAddr(addr, step string) string {
	if strings.Contains(addr, step) {
		return strings.Split(addr, step)[0]
	}
	return addr
}
