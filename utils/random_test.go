package utils

import (
	"testing"
)

func TestRandom(t *testing.T) {
	t.Log(RandInt(0, 3))
}
