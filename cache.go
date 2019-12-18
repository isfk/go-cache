package cache

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v7"
	"time"
)

// Client Client
type Client struct {
	RedisClient *redis.Client
	tags        []string
}

// RedisDriver RedisDriver
var RedisDriver *redis.Client

// NewClient NewClient
func NewClient(conf redis.Options) *Client {
	client := redis.NewClient(&conf)

	pong, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}

	fmt.Println("redis pong: ", pong)

	RedisDriver = client

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
func (c *Client) Put(key string, val interface{}, expire time.Duration) error {
	// sadd member
	for _, v := range c.tags {
		err := RedisDriver.SAdd(fmt.Sprintf("tag:%v", fmt.Sprintf("%x", md5.Sum([]byte(v)))), key).Err()
		if err != nil {
			fmt.Println("RedisDriver.SAdd err:", err)
		}
	}

	// json
	value, err := json.Marshal(val)
	if err != nil {
		fmt.Println("err:", err)
		return err
	}

	// set key:value
	err = RedisDriver.Set(key, string(value), expire).Err()
	if err != nil {
		fmt.Println("RedisDriver.Set err:", err)
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
			fmt.Println("RedisDriver.SMembers err:", err)
			return err
		}

		if len(members) > 0 {
			// delete members
			err = RedisDriver.Del(members...).Err()
			if err != nil {
				fmt.Println("RedisDriver.Del err:", err)
				return err
			}
		}

		// delete tag
		err = RedisDriver.Del(key).Err()
		if err != nil {
			fmt.Println("RedisDriver.Del err:", err)
			return err
		}
	}
	return nil
}
