package internal

import (
	"testing"
)

func TestSortRecors(t *testing.T) {
	result := []map[string]string{{"name": "test.hyuga.io", "remoteAddr": "61.155.167.50:29214", "ts": "1600335462140672515"}, {"name": "test.hyuga.io", "remoteAddr": "140.205.129.118:24129", "ts": "1600369446465646677"}, {"name": "test.hyuga.io", "remoteAddr": "140.205.129.122:34021", "ts": "1600340127877828042"}, {"name": "test.hyuga.io", "remoteAddr": "172.253.0.3:53883", "ts": "-"}, {"name": "test.hyuga.io", "remoteAddr": "140.205.129.70:42006", "ts": "1600803698367057244"}, {"name": "test.hyuga.io", "remoteAddr": "140.205.129.82:10548", "ts": "1600749576106671851"}, {"name": "test.hyuga.io", "remoteAddr": "61.155.167.41:20567", "ts": "1600745754476550746"}, {"name": "test.hyuga.io", "remoteAddr": "172.253.0.3:53883", "ts": "1600745633653942763"}}

	sorted := SortRecords(result)
	t.Log(sorted)

	for i, res := range sorted {
		t.Log(i, " ", res["ts"])
	}
}
