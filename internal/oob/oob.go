package oob

import (
	"strings"
)

func parseSid(domain, mainDomain string) string {
	i := strings.Index(domain, mainDomain)
	if i <= 0 {
		return ""
	}

	pre := strings.Split(strings.Trim(domain[:i], "."), ".")
	return pre[len(pre)-1]
}
