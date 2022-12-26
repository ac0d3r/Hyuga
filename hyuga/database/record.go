package database

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Record interface {
	Type() string
	Values() map[string]string
}

var (
	_ Record = DnsRecord{}
	_ Record = HttpRecord{}
	_ Record = JndiRecord{}
)

type DnsRecord struct {
	Name       string `json:"name"`
	RemoteAddr string `json:"remote_addr"`
	Created    int64  `json:"created"`
}

func (d DnsRecord) Values() map[string]string {
	return map[string]string{
		"name":        d.Name,
		"remote_addr": d.RemoteAddr,
		"created":     strconv.FormatInt(d.Created, 10),
	}
}

func (d DnsRecord) Type() string {
	return "dns"
}

type HttpRecord struct {
	URL        string `json:"url"`
	Method     string `json:"method"`
	RemoteAddr string `json:"remote_addr"`
	Raw        string `json:"raw"`
	Created    int64  `json:"created"`
}

func (h HttpRecord) Values() map[string]string {
	return map[string]string{
		"url":         h.URL,
		"method":      h.Method,
		"remote_addr": h.RemoteAddr,
		"raw":         h.Raw,
		"created":     strconv.FormatInt(h.Created, 10),
	}
}

func (d HttpRecord) Type() string {
	return "http"
}

type JndiRecord struct {
	Protocol   string `json:"protocol"`
	RemoteAddr string `json:"remote_addr"`
	Path       string `json:"path"`
	Created    int64  `json:"created"`
}

func (j JndiRecord) Values() map[string]string {
	return map[string]string{
		"remote_addr": j.RemoteAddr,
		"protocol":    j.Protocol,
		"path":        j.Path,
		"created":     strconv.FormatInt(j.Created, 10),
	}
}

func (l JndiRecord) Type() string {
	return "jndi"
}

func SetUserRecord(userID string, record Record, expire time.Duration) error {
	key := genRecordKey(record.Type(), userID)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*8)
	defer cancel()

	if err := defaultClient.HSet(ctx, key, record.Values()).Err(); err != nil {
		return err
	}

	return defaultClient.Expire(ctx, key, expire).Err()
}

// genRecordKey: [record-type]-[identity]-[timestamp]
func genRecordKey(t, id string) string {
	return fmt.Sprintf("%s-%s-%d", t, id, time.Now().UnixNano())
}

func DeleteRecordsByUserID(record Record, userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*8)
	defer cancel()
	allkeys, err := defaultClient.Keys(ctx, fmt.Sprintf("%s-%s-*", record.Type(), userID)).Result()
	if err != nil {
		return err
	}
	if len(allkeys) == 0 {
		return nil
	}

	return defaultClient.Del(ctx, allkeys...).Err()
}

func GetUserRecordsByUserID(record Record, userID string, filter string) ([]map[string]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*8)
	defer cancel()

	allkeys, err := defaultClient.Keys(ctx, fmt.Sprintf("%s-%s-*", record.Type(), userID)).Result()
	if err != nil {
		return nil, err
	}
	if len(allkeys) == 0 {
		return nil, nil
	}

	// sort && limit
	sort.Sort(recordKeys(allkeys))
	limit := 100
	if limit > len(allkeys) {
		limit = len(allkeys)
	}
	allkeys = allkeys[:limit]

	// filter
	sample := ""
	notfilter := len(filter) == 0
	resultList := make([]map[string]string, 0)
	for i := range allkeys {
		data, err := defaultClient.HGetAll(ctx, allkeys[i]).Result()
		if err != nil {
			return nil, err
		}
		// simply add a key
		data["key"] = strings.Split(allkeys[i], "-")[2]
		if notfilter {
			resultList = append(resultList, data)
			continue
		}
		switch record.Type() {
		case "dns":
			sample = data["name"]
		case "http":
			sample = data["url"]
		case "jdni":
			sample = data["path"]
		}
		if strings.Contains(sample, filter) {
			resultList = append(resultList, data)
		}
	}

	return resultList, nil
}

// recordKeys default descending order
type recordKeys []string

func (r recordKeys) Len() int {
	return len(r)
}

func (r recordKeys) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r recordKeys) Less(i, j int) bool {
	tsi, erri := strconv.Atoi(strings.Split(r[i], "-")[2])
	tsj, errj := strconv.Atoi(strings.Split(r[j], "-")[2])
	if erri != nil && errj == nil {
		return false
	}
	if erri == nil && errj != nil {
		return true
	}
	return tsi > tsj
}
