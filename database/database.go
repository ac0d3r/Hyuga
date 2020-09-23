package database

import (
	"Hyuga/conf"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/gommon/log"
)

type record struct {
	rdb *redis.Client
}

var (
	ctx = context.Background()
	// Recorder Hyuga redis recorder
	Recorder *record = newRecorder()
	// DefaultLimit `GetRecords` default limit
	DefaultLimit int = 50
)

func logformat(msg string) string {
	return fmt.Sprintf("Recorder redis %s:%d %s", conf.RedisHost, conf.RedisPort, msg)
}

func filterRecordType(rtype string) error {
	var err error
	switch rtype {
	case "dns":
		err = nil
	case "http":
		err = nil
	default:
		err = fmt.Errorf("Unsupported record type: '%s'", rtype)
	}
	return err
}

func genRecordsKey(rtype, identity string) string {
	return fmt.Sprintf("%s-%s-%d", rtype, identity, time.Now().UnixNano())
}

func newRecorder() *record {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.RedisHost, conf.RedisPort),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Error(logformat(err.Error()))
	}
	log.Debug(logformat(pong))
	return &record{rdb: rdb}
}

func (rc *record) IdentityExist(identity string) bool {
	values, err := rc.rdb.Keys(ctx, fmt.Sprintf("user-%s-*", identity)).Result()
	if len(values) == 0 {
		return false
	}
	if err != nil {
		log.Error(logformat(err.Error()))
		return false
	}
	return true
}

func (rc *record) UserExist(identity, token string) bool {
	_, err := rc.rdb.Get(ctx, fmt.Sprintf("user-%s-%s", identity, token)).Result()
	if err == redis.Nil {
		return false
	} else if err != nil {
		log.Error(logformat(err.Error()))
		return false
	}
	return true
}

func (rc *record) CreateUser(identity, token string) error {
	err := rc.rdb.SetNX(ctx, fmt.Sprintf("user-%s-%s", identity, token), time.Now().UnixNano(), 0).Err()
	if err != nil {
		log.Error(logformat(err.Error()))
	}
	return err
}

func (rc *record) Record(rtype, identity string, data map[string]interface{}) error {
	// do not record when identity does not exist
	if !rc.IdentityExist(identity) {
		return fmt.Errorf("Not Found user identity: '%s'", identity)
	}
	if err := filterRecordType(rtype); err != nil {
		return err
	}
	switch rtype {
	case "dns":
		return rc.recordDNS(identity, data)
	case "http":
		return rc.recordHTTP(identity, data)
	}
	return nil
}

func (rc *record) recordDNS(identity string, data map[string]interface{}) error {
	key := genRecordsKey("dns", identity)
	err := rc.rdb.HMSet(ctx, key, data).Err()
	if err != nil {
		log.Error(logformat(err.Error()))
	} else {
		rc.rdb.Expire(ctx, key, conf.RecordExpiration)
	}
	return err
}

func (rc *record) recordHTTP(identity string, data map[string]interface{}) error {
	key := genRecordsKey("http", identity)
	err := rc.rdb.HMSet(ctx, key, data).Err()
	if err != nil {
		log.Error(logformat(err.Error()))
	} else {
		rc.rdb.Expire(ctx, key, conf.RecordExpiration)
	}
	return err
}

func (rc *record) GetRecords(rtype, token, filter string) ([]map[string]string, error) {
	if err := filterRecordType(rtype); err != nil {
		return []map[string]string{}, err
	}
	// get user identity with `token`
	users, err := rc.rdb.Keys(ctx, fmt.Sprintf("user-*-%s", token)).Result()
	if err != nil || len(users) == 0 {
		if err != nil {
			log.Error(logformat(err.Error()))
		}
		return []map[string]string{}, err
	}
	identity := strings.Split(users[0], "-")[1]
	// get user records when type is `rtype`
	recordKeys, err := rc.rdb.Keys(ctx, fmt.Sprintf("%s-%s-*", rtype, identity)).Result()
	if err != nil || len(recordKeys) == 0 {
		if err != nil {
			log.Error(logformat(err.Error()))
		}
		return []map[string]string{}, err
	}
	log.Debug("GetRecords", recordKeys)

	var result []map[string]string
	// set maximum record length limit
	limit := DefaultLimit
	if len(recordKeys) < DefaultLimit {
		limit = len(recordKeys)
	}
	for _, key := range recordKeys[:limit] {
		ts := strings.Split(key, "-")[2]
		data, err := rc.rdb.HGetAll(ctx, key).Result()
		if err != nil {
			log.Error(logformat(err.Error()))
			continue
		}
		data["ts"] = ts
		log.Debug(data)
		// filter records
		if filter == "" {
			result = append(result, data)
			continue
		}
		switch rtype {
		case "dns":
			if strings.Contains(data["name"], filter) {
				result = append(result, data)
			}
		case "http":
			if strings.Contains(data["url"], filter) {
				result = append(result, data)
			}
		}
	}
	return result, nil
}

func (rc *record) Close() {
	rc.rdb.Close()
}
