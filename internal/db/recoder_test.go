package db

import (
	"fmt"
	"testing"
)

func TestRecod(t *testing.T) {
	r := NewRecorder()
	for i := 0; i < 66; i++ {
		r.Record("test", fmt.Sprintf("%s-%d", "test", i))
	}

	for i := 0; i < 77; i++ {
		r.Record("admin", fmt.Sprintf("%s-%d", "admin", i))
	}

	t.Log(r.Get("test"))
	t.Log(r.Get("admin"))
}
