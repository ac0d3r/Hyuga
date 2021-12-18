package config

import (
	"testing"
)

func TestSetFromFile(t *testing.T) {
	err := SetFromYaml("config.yaml")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(DebugMode)
	t.Log(RedisDsn)
	t.Log(RecordExpiration)
	t.Log(MainDomain)
	t.Log(DefaultIP)
	t.Log(NSDomain)
	t.Log(DefaultIP)
}
