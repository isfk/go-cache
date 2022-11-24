package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// Cache Cache
type Cache struct {
	RedisClient        *redis.Client
	RedisClusterClient *redis.ClusterClient
	tags               []string
}

// RedisDriver RedisDriver
var RedisDriver *redis.Client

// RedisClusterDriver RedisClusterDriver
var RedisClusterDriver *redis.ClusterClient

// New New
func New(ctx context.Context, client *redis.Client) *Cache {
	pong, err := client.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}

	fmt.Println("redis pong: ", pong)

	RedisDriver = client

	return &Cache{
		RedisClient: RedisDriver,
	}
}

// NewClusterClient NewClusterClient
func NewClusterClient(ctx context.Context, client *redis.ClusterClient) *Cache {
	pong, err := client.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}

	fmt.Println("redis pong: ", pong)

	RedisClusterDriver = client

	return &Cache{
		RedisClusterClient: RedisClusterDriver,
	}
}

// Tag .Tag()
func (c *Cache) Tag(tag ...string) *Cache {
	c.tags = tag
	return c
}

// Set .Tag().Set()
func (c *Cache) Set(ctx context.Context, key string, val interface{}, expire time.Duration) error {
	for _, v := range c.tags {
		err := RedisDriver.SAdd(ctx, v, key).Err()
		if err != nil {
			fmt.Println("RedisDriver.SAdd err:", err)
		}
	}

	value, err := json.Marshal(val)
	if err != nil {
		fmt.Println("err:", err)
		return err
	}

	err = RedisDriver.Set(ctx, key, string(value), expire).Err()
	if err != nil {
		fmt.Println("RedisDriver.Set err:", err)
		return err
	}

	return nil
}

// Get .Get()
func (c *Cache) Get(ctx context.Context, key string, val interface{}) error {
	jsonStr, err := RedisDriver.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil
		}
		fmt.Println("RedisDriver.Get err:", err)
		return err
	}
	err = json.Unmarshal([]byte(jsonStr), &val)

	if err != nil {
		fmt.Println("json.Unmarshal err:", err)
		return err
	}

	return nil
}

// Flush .Tag().Flush()
func (c *Cache) Flush(ctx context.Context) error {
	for _, key := range c.tags {
		members, err := RedisDriver.SMembers(ctx, key).Result()
		if err != nil {
			fmt.Println("RedisDriver.SMembers err:", err)
			return err
		}

		if len(members) > 0 {
			err = RedisDriver.Del(ctx, members...).Err()
			if err != nil {
				fmt.Println("RedisDriver.Del err:", err)
				return err
			}
		}

		err = RedisDriver.Del(ctx, key).Err()
		if err != nil {
			fmt.Println("RedisDriver.Del err:", err)
			return err
		}
	}
	return nil
}
