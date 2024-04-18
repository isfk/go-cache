package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Options struct {
	prefix  string
	expired time.Duration
}

type Option func(*Options)

// Cache Cache
type Cache struct {
	redis        *redis.Client
	redisCluster *redis.ClusterClient
	tags         []string

	prefix  string
	expired time.Duration
}

func WithPrefix(prefix string) Option {
	return func(o *Options) {
		o.prefix = prefix
	}
}

func WithExpired(exp time.Duration) Option {
	return func(o *Options) {
		o.expired = exp
	}
}

// New New
func New(ctx context.Context, client *redis.Client, opts ...Option) *Cache {
	pong, err := client.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}

	fmt.Println("redis pong: ", pong)

	o := &Options{}
	for _, opt := range opts {
		opt(o)
	}

	return &Cache{
		redis:   client,
		prefix:  o.prefix,
		expired: o.expired,
	}
}

// NewCluster NewCluster
func NewCluster(ctx context.Context, client *redis.ClusterClient) *Cache {
	pong, err := client.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}

	fmt.Println("redis pong: ", pong)

	return &Cache{
		redisCluster: client,
	}
}

// Tag .Tag()
func (c *Cache) Tag(tag ...string) *Cache {
	c.tags = tag
	return c
}

// Tag .AddTag()
func (c *Cache) AddTag(tag ...string) *Cache {
	c.tags = append(c.tags, tag...)
	return c
}

// Set .Tag().Set()
func (c *Cache) Set(ctx context.Context, key string, val interface{}) error {
	if len(c.prefix) > 0 {
		key = c.prefix + ":" + key
	}

	_, err := c.redis.TxPipelined(ctx, func(p redis.Pipeliner) error {
		for _, v := range c.tags {
			err := p.SAdd(ctx, c.prefix+":"+v, key).Err()
			if err != nil {
				fmt.Println(fmt.Errorf("p.SAdd err %v", err))
				return err
			}
		}

		value, err := json.Marshal(val)
		if err != nil {
			fmt.Println("json.Marshal err:", err)
			return err
		}

		err = p.Set(ctx, key, string(value), c.expired).Err()
		if err != nil {
			fmt.Println("p.Set err:", err)
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// Get .Get()
func (c *Cache) Get(ctx context.Context, key string, val interface{}) error {
	if len(c.prefix) > 0 {
		key = c.prefix + ":" + key
	}

	jsonStr, err := c.redis.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil
		}
		fmt.Println("c.redis.Get err:", err)
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
	_, err := c.redis.TxPipelined(ctx, func(p redis.Pipeliner) error {
		for _, v := range c.tags {
			members, err := c.redis.SMembers(ctx, c.prefix+":"+v).Result()
			if err != nil {
				fmt.Println("c.redis.SMembers err:", err)
				return err
			}

			if len(members) > 0 {
				err = p.Del(ctx, members...).Err()
				if err != nil {
					fmt.Println("p.Del err:", err)
					return err
				}
			}

			err = p.Del(ctx, c.prefix+":"+v).Err()
			if err != nil {
				fmt.Println("p.Del err:", err)
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
