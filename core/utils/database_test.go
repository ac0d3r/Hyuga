package utils

import (
	"testing"
)

func TestRecorder(t *testing.T) {
	defer Recorder.Close()
	var (
		err      error
		identity = "test"
		token    = "testtoken"
		records  []map[string]string
	)

	Recorder.CreateUser(identity, token)

	t.Log("IdentityExist ", identity, Recorder.IdentityExist(identity))
	t.Log("UserExist ", identity, token, Recorder.UserExist(identity, token))

	dnsData := map[string]interface{}{"name": "test.hyuga.io", "remote_addr": "127.0.0.1:1314"}
	Recorder.Record("dns", "test", dnsData)

	httpData := map[string]interface{}{"url": "http://test.hyuga.io", "remote_addr": "127.0.0.1:1314", "method": "GET", "cookies": "a=123; b=456"}
	Recorder.Record("http", identity, httpData)

	err = Recorder.Record("others", identity, dnsData)
	t.Log("Record others", err)

	records, err = Recorder.GetRecords("dns", token, "Not Found")
	if err != nil {
		t.Log("records error:", err)
	}
	t.Log("records:", records)
}
