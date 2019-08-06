package cache

import (
	"fmt"

	"github.com/go-redis/redis"
)

// RedisDriver RedisDriver
var RedisDriver *redis.Client

// InitRedis InitRedis
func InitRedis() {
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "sdfsdf", // no password set
		DB:       0,        // use default DB
	})

	pong, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("redis pong: ", pong)
	RedisDriver = client
}
