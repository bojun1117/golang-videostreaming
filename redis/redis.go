package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var (
	rdb *redis.Client
)

const (
	// Initialize connection constants.
	ADDR     = "videostreaming.6iwt1k.ng.0001.apne1.cache.amazonaws.com:6379"
	PASSWORD = ""
	DATABASE = 0
)

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     ADDR,
		Password: PASSWORD,
		DB:       DATABASE,
	})
}

func Check(key string) string {
	ctx := context.Background()
	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		return "keynull"
	}
	return val
}

func Setkey(key string, val string) error {
	ctx := context.Background()
	err := rdb.Set(ctx, key, val, 0).Err()
	if err != nil {
		return err
	}
	return nil
}
