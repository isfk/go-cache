package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// Client Client
type Client struct {
	RedisClient        *redis.Client
	RedisClusterClient *redis.ClusterClient
	tags               []string
}

// RedisDriver RedisDriver
var RedisDriver *redis.Client

// RedisClusterDriver RedisClusterDriver
var RedisClusterDriver *redis.ClusterClient

// NewClient NewClient
func NewClient(ctx context.Context, client *redis.Client) *Client {
	pong, err := client.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}

	fmt.Println("redis pong: ", pong)

	RedisDriver = client

	return &Client{
		RedisClient: RedisDriver,
	}
}

// NewClusterClient NewClusterClient
func NewClusterClient(ctx context.Context, client *redis.ClusterClient) *Client {
	pong, err := client.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}

	fmt.Println("redis pong: ", pong)

	RedisClusterDriver = client

	return &Client{
		RedisClusterClient: RedisClusterDriver,
	}
}

// Tag .Tag()
func (c *Client) Tag(tag ...string) *Client {
	c.tags = tag
	return c
}

// Set .Tag().Set()
func (c *Client) Set(ctx context.Context, key string, val interface{}, expire time.Duration) error {
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
func (c *Client) Get(ctx context.Context, key string, val interface{}) error {
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
func (c *Client) Flush(ctx context.Context) error {
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
