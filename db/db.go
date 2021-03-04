package db

import (
	"context"
	"fmt"
	"sync"
	"time"

	"hyuga/conf"

	"github.com/go-redis/redis/v8"
)

// Limit 查询记录限制
const Limit = 50

const (
	// DNSType dns
	DNSType = "dns"
	// HTTPType http
	HTTPType = "http"
)

var (
	rc   *redis.Client
	once sync.Once
)

// New 创建 redis client
func New() {
	once.Do(func() {
		rc = redis.NewClient(&redis.Options{
			Addr:     conf.RedisAddr,
			Password: "", // no password set
			DB:       0,  // use default DB
		})
	})
}

func logformat(msg string) string {
	return fmt.Sprintf("redis %s %s", conf.RedisAddr, msg)
}

func genRecordKey(rtype, identity string) string {
	// RecordsKey: [record-type]-[identity]-[timestamp]
	return fmt.Sprintf("%s-%s-%d", rtype, identity, time.Now().UnixNano())
}

func recordDNS(ctx context.Context, identity string, data map[string]interface{}) error {
	// +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
	// {
	//	"name": [dns query name],
	//	"remoteAddr": [remote address]
	// }
	// +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+

	key := genRecordKey(DNSType, identity)
	err := rc.HMSet(ctx, key, data).Err()
	if err != nil {
		return err
	}
	return rc.Expire(ctx, key, conf.RecordExpiration).Err()
}

func recordHTTP(ctx context.Context, identity string, data map[string]interface{}) error {
	// record other requests from `*.huyga.io`
	// +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
	// {
	// 	"url":        [request uri],
	// 	"method":     [request method],
	// 	"remoteAddr": [request remote address],
	// 	"cookies":    [request cookies],
	// }
	// +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
	key := genRecordKey(HTTPType, identity)
	err := rc.HMSet(ctx, key, data).Err()
	if err != nil {
		return err
	}
	return rc.Expire(ctx, key, conf.RecordExpiration).Err()
}
