package config

import (
	"testing"
)

func TestSetFromFile(t *testing.T) {
	err := SetFromYaml("config.yaml")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v \n", C)
}
