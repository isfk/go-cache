package cache

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

// Client Client
type Client struct {
	baseClient *redis.Client
	tags       []string
}

// RedisDriver RedisDriver
var RedisDriver *redis.Client

// NewClient NewClient
func NewClient() *Client {
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
	for _, v := range c.tags {
		err := RedisDriver.SAdd(fmt.Sprintf("tag:%v", fmt.Sprintf("%x", md5.Sum([]byte(v)))), key).Err()
		if err != nil {
			fmt.Println("err:", err)
		}
	}

	value, err := json.Marshal(val)
	if err != nil {
		fmt.Println("err:", err)
		return err
	}

	err = RedisDriver.Set(fmt.Sprintf("key:%v", key), string(value), time.Hour).Err()

	if err != nil {
		fmt.Println("err:", err)
		return err
	}

	return nil
}

// Clear .Tag().Clear()
func (c *Client) Clear() error {
	for _, val := range c.tags {
		// get key
		key := fmt.Sprintf("tag:%v", fmt.Sprintf("%x", md5.Sum([]byte(val))))

		// get members
		members, err := RedisDriver.SMembers(key).Result()
		if err != nil {
			fmt.Println("err:", err)
			return err
		}

		// delete members
		err = RedisDriver.Del(members...).Err()
		if err != nil {
			fmt.Println("err:", err)
			return err
		}

		// delete tag
		err = RedisDriver.Del(key).Err()
		if err != nil {
			fmt.Println("err:", err)
			return err
		}
	}
	return nil
}
