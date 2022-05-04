package caching

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/mock"
)

func NewRedis(host string, pass string) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: pass,
		DB:       0,
	})

	return rdb
}

type MockRedis struct {
	mock.Mock
}

func (mr *MockRedis) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	args := mr.Called(ctx, key, value, expiration)
	return args.Get(0).(*redis.StatusCmd)
}

func (mr *MockRedis) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	args := mr.Called(ctx, keys)
	return args.Get(0).(*redis.IntCmd)
}

func (mr *MockRedis) Get(ctx context.Context, key string) *redis.StringCmd {
	args := mr.Called(ctx, key)
	return args.Get(0).(*redis.StringCmd)
}

func (mr *MockRedis) Keys(ctx context.Context, pattern string) *redis.StringSliceCmd {
	args := mr.Called(ctx, pattern)
	return args.Get(0).(*redis.StringSliceCmd)
}

// func (mr *MockRedis) Result() (string, error) {
// 	args := mr.Called()
// 	return args.Get(0).(string), args.Error(1)
// }
