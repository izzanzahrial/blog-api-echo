package redis

import (
	"github.com/go-redis/redis/v8"
)

// func NewRedis() {
// 	rdb := redis.NewClient(&redis.Options{
// 		Addr: "localhost:6379",
// 	})

// 	mycache := cache.New(&cache.Options{
// 		Redis: rdb,
// 	})

// 	obj := 1
// 	err := mycache.Once(&cache.Item{
// 		Key: "key",
// 		Value: obj,
// 		Do: func(i *cache.Item) (interface{}, error) {
// 			return obj, nil
// 		},
// 	})

// 	if err != nil {
// 		blabla
// 	}
// }

func NewRedis(host string, pass string) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: pass,
		DB:       0,
	})

	return rdb
}
