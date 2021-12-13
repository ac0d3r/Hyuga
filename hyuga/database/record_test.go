package database

import (
	"fmt"
	"testing"
	"time"
)

func Test_setRecord(t *testing.T) {
	for i := 0; i < 15; i++ {
		record := DnsRecord{
			Name:       fmt.Sprintf("%d.hyuga.io", i),
			RemoteAddr: "127.0.0.1",
		}
		t.Log(SetUserRecord("123", record, time.Hour))
	}
}

func Test_getRecord(t *testing.T) {
	t.Log(GetUserRecordsByUserID(DnsRecord{}, "123", ""))
}
