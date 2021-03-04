package random

import (
	"math/rand"
	"time"
)

const (
	// Charset lower case character set
	Charset = "abcdefghijklmnopqrstuvwxyz"
	// UpperCharset upper case character set
	UpperCharset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	// Number 0-9 number set
	Number = "0123456789"
)

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

// StringWithCharset  random string
func StringWithCharset(length int, charset string) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// String random string
func String(length int) string {
	c := Charset + UpperCharset + Number
	return StringWithCharset(length, c)
}

// Int random integer between min~max
func Int(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}
