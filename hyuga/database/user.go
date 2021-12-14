package database

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
)

type User struct {
	ID        string   `json:"id"`
	Token     string   `json:"token"`
	IPs       []string `json:"ip"`
	Created   int64    `json:"created"`
	RDnsTimes int64    `json:"rdnstimes"` // rebinding dns query of times
}

// TODO use reflect
func (u *User) values() (map[string]string, error) {
	values := map[string]string{
		"id":        u.ID,
		"token":     u.Token,
		"created":   strconv.FormatInt(u.Created, 10),
		"ip":        "[]",
		"rdnstimes": "0",
	}
	if len(u.IPs) > 0 {
		ips, err := json.Marshal(u.IPs)
		if err == nil {
			return nil, err
		}
		values["ip"] = string(ips)
	}
	return values, nil
}

// TODO use reflect
func (u *User) makeWithValues(values map[string]string) error {
	for k, v := range values {
		switch k {
		case "id":
			u.ID = v
		case "token":
			u.Token = v
		case "created":
			c, err := strconv.Atoi(v)
			if err == nil {
				u.Created = int64(c)
			}
		case "ip":
			ip := make([]string, 0)
			if err := json.Unmarshal([]byte(v), &ip); err != nil {
				return err
			}
			u.IPs = ip
		}
	}
	return nil
}

func genUserKey(id string) string {
	return "user-" + id
}

func (u *User) Key() string {
	return genUserKey(u.ID)
}

func CreateUser(user *User) error {
	if UserExist(user.ID) {
		return ErrUserAlreadyExists
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// set user token index
	if err := defaultClient.Set(ctx, user.Token, user.Key(), 0).Err(); err != nil {
		return err
	}
	values, err := user.values()
	if err != nil {
		return err
	}
	return defaultClient.HSet(ctx, user.Key(), values).Err()
}

func DeleteUserByUserID(userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	r := defaultClient.HGet(ctx, genUserKey(userID), "token")
	token, err := r.Result()
	if err != nil {
		return err
	}

	return defaultClient.Del(ctx, genUserKey(userID), token).Err()
}

func GetUserByToken(token string) (*User, error) {
	var user = &User{}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	key, err := defaultClient.Get(ctx, token).Result()
	if err != nil {
		return nil, err
	}

	data, err := defaultClient.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	if err := user.makeWithValues(data); err != nil {
		return nil, err
	}
	return user, nil
}

func UserExist(userID string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := defaultClient.HGet(ctx, genUserKey(userID), "id").Result()
	return err == nil
}

func SetUserDNSRebinding(userID string, ip []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	var ipstr string
	if len(ip) == 0 {
		ipstr = "[]"
	} else {
		ips, err := json.Marshal(ip)
		if err != nil {
			return err
		}
		ipstr = string(ips)
	}
	return defaultClient.HSet(ctx, genUserKey(userID), "ip", ipstr).Err()
}

func GetUserDNSRebinding(userID string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	bytes, err := defaultClient.HGet(ctx, genUserKey(userID), "ip").Bytes()
	if err != nil {
		return nil, err
	}

	if len(bytes) == 0 {
		return []string{}, nil
	}

	ip := make([]string, 0)
	if err := json.Unmarshal(bytes, &ip); err != nil {
		return nil, err
	}
	return ip, err
}

func SetUserDnsRebindingTimes(userID string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	t, err := defaultClient.HGet(ctx, genUserKey(userID), "rdnstimes").Int64()
	if err != nil {
		return 0, err
	}

	return t, defaultClient.HSet(ctx, genUserKey(userID), "rdnstimes", strconv.FormatInt(t+1, 10)).Err()
}
