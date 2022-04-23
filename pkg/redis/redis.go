package redis

import (
	"github.com/go-redis/redis/v8"
)

func NewRedis(host string, pass string) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: pass,
		DB:       0,
	})

	return rdb
}
