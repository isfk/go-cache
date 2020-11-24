package cache

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

// Client Client
type Client struct {
	RedisClient *redis.Ring
	tags        []string
}

// RedisDriver RedisDriver
var RedisDriver *redis.Ring

// NewClient NewClient
func NewClient(ctx context.Context, red *redis.Ring) *Client {
	pong, err := red.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}

	fmt.Println("redis pong: ", pong)

	RedisDriver = red

	return &Client{
		RedisClient: RedisDriver,
	}
}

// Tag .Tag()
func (c *Client) Tag(tag ...string) *Client {
	c.tags = tag
	return c
}

// Put .Tag().Put()
func (c *Client) Put(ctx context.Context, key string, val interface{}, expire time.Duration) error {
	for _, v := range c.tags {
		err := RedisDriver.SAdd(ctx, fmt.Sprintf("tag:%v", fmt.Sprintf("%x", md5.Sum([]byte(v)))), key).Err()
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

// Clear .Tag().Get()
func (c *Client) Get(ctx context.Context, key string, val interface{}) (interface{}, error) {
	jsonStr, err := RedisDriver.Get(ctx, key).Result()
	if err != nil {
		fmt.Println("RedisDriver.Get err:", err)
		return nil, err
	}
	err = json.Unmarshal([]byte(jsonStr), &val)

	if err != nil {
		fmt.Println("json.Unmarshal err:", err)
		return nil, err
	}

	return val, nil
}

// Clear .Tag().Clear()
func (c *Client) Clear(ctx context.Context) error {
	for _, val := range c.tags {
		key := fmt.Sprintf("tag:%v", fmt.Sprintf("%x", md5.Sum([]byte(val))))

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
