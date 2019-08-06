package cache

import (
	"crypto/md5"
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

// Client Client
type Client struct {
	baseClient *redis.Client
	tags       []string
}

// NewClient NewClient
func NewClient() *Client {
	InitRedis()

	return &Client{
		baseClient: RedisDriver,
	}
}

// Tag .Tag()
func (c *Client) Tag(tag ...string) *Client {
	c.tags = tag
	return c
}

// Put .Tag().Put()
func (c *Client) Put(key string, val interface{}, expire int64) error {
	InitRedis()

	for _, v := range c.tags {
		err := RedisDriver.SAdd(key, fmt.Sprintf("tag:%v", fmt.Sprintf("%x", md5.Sum([]byte(v))))).Err()

		if err != nil {
			fmt.Println("err:", err)
		}
	}

	err := RedisDriver.Set(key, val, time.Hour).Err()

	if err != nil {
		fmt.Println("err:", err)
	}

	return nil
}

// Clear .Tag().Clear()
func (c *Client) Clear(key string) error {
	InitRedis()

	err := RedisDriver.Del(key).Err()

	if err != nil {
		fmt.Println("err:", err)
	}
	return nil
}
