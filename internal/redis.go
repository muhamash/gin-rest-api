package redisclient

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()

func NewClient(redisURL string) *redis.Client {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		panic(err)
	}
	return redis.NewClient(opt)
}