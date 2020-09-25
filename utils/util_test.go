package utils

import (
	"testing"
)

func TestParseRemoteAddr(t *testing.T) {
	t.Log(ParseRemoteAddr("127.0.0.1:5000", ":"))
}
