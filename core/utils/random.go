package utils

import (
	"math/rand"
	"time"
)

const (
	// Charset lower case character set
	Charset = "abcdefghijklmnopqrstuvwxyz"
	// UpperCharset upper case character set
	UpperCharset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	// Number 0-9 number set
	Number = "0123456789"
)

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

// RandomStringWithCharset  random string
func RandomStringWithCharset(length int, charset string) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// RandomString random string
func RandomString(length int) string {
	c := Charset + UpperCharset + Number
	return RandomStringWithCharset(length, c)
}

// RandInt random integer between min~manx
func RandInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}
