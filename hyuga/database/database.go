package database

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	defaultClient *redis.Client
)

func Init(dsn string) error {
	var (
		opt *redis.Options
		err error
	)
	// compatible with docker link
	if dsn == "redis:6379" {
		opt = &redis.Options{Addr: dsn, Password: "", DB: 0}
	} else {
		if opt, err = redis.ParseURL(dsn); err != nil {
			return err
		}
	}
	defaultClient = redis.NewClient(opt)
	return ping()
}

func ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if _, err := defaultClient.Ping(ctx).Result(); err != nil {
		return err
	}

	return nil
}
