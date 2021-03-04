package random

import (
	"testing"
)

func TestRandom(t *testing.T) {
	t.Log(StringWithCharset(3, Charset))
	t.Log(StringWithCharset(4, UpperCharset))
	t.Log(StringWithCharset(5, Number))

	t.Log(String(10))

	t.Log(Int(0, 3))
}
