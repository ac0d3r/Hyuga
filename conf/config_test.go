package conf

import "testing"

func TestSetFromFile(t *testing.T) {
	t.Log(SetFromFile("config-test.yml"))

	t.Log(AppEnv)
	t.Log(RedisAddr)
	t.Logf("%#v", Domain)
	t.Log(RecordExpiration)
}
