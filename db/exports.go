package db

import (
	"context"
	"fmt"
	"strings"
	"time"

	"hyuga/internal/random"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/gommon/log"
)

var (
	// ErrUnsupportedRecordType 不支持的记录类型
	ErrUnsupportedRecordType = fmt.Errorf("Unsupported record type")
	// ErrUserNotFound user不存在
	ErrUserNotFound = fmt.Errorf("user not found")
)

// Ping test redis
func Ping(ctx context.Context) (string, error) {
	pong, err := rc.Ping(ctx).Result()
	if err != nil {
		return "", err
	}
	return pong, nil
}

// CreateUser 创建用户
func CreateUser(ctx context.Context, identity, token string) error {
	// "user-[identity]-[token]": timestamp
	return rc.SetNX(ctx, fmt.Sprintf("user-%s-%s", identity, token), time.Now().UnixNano(), 0).Err()
}

// UserExist 判断用户是否存在
func UserExist(ctx context.Context, identity, token string) (bool, error) {
	_, err := rc.Get(ctx, fmt.Sprintf("user-%s-%s", identity, token)).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// IdentityExist 判断 user id 是否存在
func IdentityExist(ctx context.Context, identity string) (bool, error) {
	values, err := rc.Keys(ctx, fmt.Sprintf("user-%s-*", identity)).Result()
	if err != nil {
		return false, err
	}
	if len(values) == 0 {
		return false, nil
	}
	return true, nil
}

// SetUserDNSRebinding 设置用户 dns-rebinding ips
func SetUserDNSRebinding(ctx context.Context, identity, token string, ips []interface{}) error {
	var (
		err  error
		llen int64
	)
	exist, err := UserExist(ctx, identity, token)
	if err != nil {
		return err
	}
	if !exist {
		return ErrUserNotFound
	}

	if llen, err = rc.LLen(ctx, identity).Result(); err != nil {
		return err
	}
	// remove all dns rebinding
	for i := 0; i < int(llen); i++ {
		if _, err = rc.LPop(ctx, identity).Result(); err != nil {
			return err
		}
	}
	if len(ips) >= 1 {
		if _, err = rc.LPush(ctx, identity, ips...).Result(); err != nil {
			return err
		}
	}
	return nil
}

// GetUserDNSRebinding 获取用户 dns-rebinding ips
// @identity 用户id
// @rand 是否随机打乱
func GetUserDNSRebinding(ctx context.Context, identity string, rand bool) (ips []string, err error) {
	exist, err := IdentityExist(ctx, identity)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, ErrUserNotFound
	}

	llen, err := rc.LLen(ctx, identity).Result()
	if err != nil {
		return nil, err
	}
	if llen == 0 {
		return ips, nil
	}

	if ips, err = rc.LRange(ctx, identity, 0, llen).Result(); err != nil {
		return nil, err
	}
	if !rand {
		return ips, err
	}
	// random get dns-rebinding ip
	var index int64
	if llen == 1 {
		index = 0
	} else {
		index = int64(random.Int(0, int(llen-1)))
	}
	return []string{ips[index]}, err
}

// SetRecord http 请求 or dns 查询
func SetRecord(ctx context.Context, rtype, identity string, data map[string]interface{}) error {
	var (
		err error
	)
	exist, err := IdentityExist(ctx, identity)
	if err != nil {
		return err
	}
	if !exist {
		return ErrUserNotFound
	}

	switch rtype {
	case DNSType:
		return recordDNS(ctx, identity, data)
	case HTTPType:
		return recordHTTP(ctx, identity, data)
	default:
		return ErrUnsupportedRecordType
	}
}

// GetRecord 查询记录
func GetRecord(ctx context.Context, rtype, token, filter string) ([]map[string]string, error) {
	userKeys, err := rc.Keys(ctx, "user-*-"+token).Result()
	if err != nil {
		return nil, err
	}
	if len(userKeys) == 0 {
		return nil, nil
	}

	identity := strings.Split(userKeys[0], "-")[1]
	recordKeys, err := rc.Keys(ctx, fmt.Sprintf("%s-%s-*", rtype, identity)).Result()
	if err != nil {
		return nil, err
	}
	if len(recordKeys) == 0 {
		return nil, nil
	}

	var result []map[string]string
	// set maximum record length limit
	limit := Limit
	if len(recordKeys) < limit {
		limit = len(recordKeys)
	}
	for _, key := range recordKeys[:limit] {
		data, err := rc.HGetAll(ctx, key).Result()
		if err != nil {
			return nil, err
		}

		data["ts"] = strings.Split(key, "-")[2]
		log.Debug(data)

		// filter record
		if filter == "" {
			result = append(result, data)
			continue
		}
		switch rtype {
		case DNSType:
			if strings.Contains(data["name"], filter) {
				result = append(result, data)
			}
		case HTTPType:
			if strings.Contains(data["url"], filter) {
				result = append(result, data)
			}
		}
	}
	return result, nil
}

// Close 关闭连接
func Close() {
	rc.Close()
}
