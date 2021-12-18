package util

import (
	"math/rand"
	"time"
)

var (
	digitSeeds       = []byte("0123456789")
	upperLetterSeeds = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	lowerLetterSeeds = []byte("abcdefghijklmnopqrstuvwxyz")
)

func RandomDigitString(length int) string {
	return randomString(digitSeeds, length)
}

func RandomID(length int) string {
	return randomString(append(digitSeeds, lowerLetterSeeds...), length)
}

func RandomString(length int) string {
	return randomString(append(append(digitSeeds, upperLetterSeeds...), lowerLetterSeeds...), length)
}

func RandomInt(min, max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(max-min+1) + min
}

func randomString(seeds []byte, length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	slice := make([]byte, 0)
	for i := 0; i < length; i++ {
		slice = append(slice, seeds[r.Intn(len(seeds))])
	}
	return string(slice)
}
