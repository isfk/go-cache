package cache

import (
	"crypto/md5"
	"fmt"
	"github.com/go-redis/cache"
	"time"
)

// Client Client
type Client struct {
	baseClient *cache.Codec
	tags   []string
}

// NewClient NewClient
func NewClient() *Client {
	GetCodec()

	return &Client{
		baseClient: Driver,
	}
}

// Handler .Handler()
func (c *Client) Handler() *cache.Codec {
	GetCodec()
	return Driver
}

// Tag .Tag()
func (c *Client) Tag(tag ...string) *Client {
	c.tags = tag
	return c
}

// Put .Tag().Put()
func (c *Client) Put(key string, val interface{}, expire int64) error {
	GetCodec()

	for _, v := range c.tags {
		err := Driver.Set(&cache.Item{
			Key:        key,
			Object:     fmt.Sprintf("tag:%v", fmt.Sprintf("%x", md5.Sum([]byte(v)))),
			Expiration: time.Hour,
		})

		if err != nil {
			fmt.Println("err:", err)
		}
	}

	err := Driver.Set(&cache.Item{
		Key:        key,
		Object:     val,
		Expiration: time.Hour,
	})

	if err != nil {
		fmt.Println("err:", err)
	}

	return nil
}

// Clear .Tag().Clear()
func (c *Client) Clear(key string) error {
	GetCodec()

	err := Driver.Delete(key)

	if err != nil {
		return err
	}

	return nil
}
