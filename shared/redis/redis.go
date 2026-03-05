package redisdb

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var Ctx = context.Background()
var Client *redis.Client

func ConnectRedis() {

	Client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

}
