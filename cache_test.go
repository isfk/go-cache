package cache

import (
	"context"
	"testing"
	"time"

	redisV9 "github.com/go-redis/redis/v9"
)

func TestCache(t *testing.T) {
	rdb := redisV9.NewClient(&redisV9.Options{
		Addr: "127.0.0.1:6379",
	})

	key := "key1"
	c := New(context.Background(), rdb, WithPrefix("test"), WithExpired(3600*time.Second))
	c.redis.Set(context.Background(), key, "value", c.expired)
	v, err := c.redis.Get(context.Background(), key).Result()
	if err != nil {
		t.Error(err)
	}
	if v != "value" {
		t.Errorf("expected value to be 'value', got '%s'", v)
	}
	c.redis.Del(context.Background(), key)

	_, err = c.redis.Get(context.Background(), key).Result()
	if err != nil {
		if err != redisV9.Nil {
			t.Error(err)
		}
	}

	// tag
	key = "key:user:1"
	tags := []string{
		"tag:all",
		"tag:user:1",
	}

	type Data struct {
		Name string
		Age  int64
	}

	data1 := &Data{Name: "jack", Age: 18}
	err = c.Tag([]string{
		"tag:all",
		"tag:user:1",
	}...).Set(context.Background(), key, data1)
	if err != nil {
		t.Error(err)
	}

	data2 := &Data{}
	err = c.Get(context.Background(), key, data2)
	if err != nil {
		t.Error(err)
	}
	if data2.Name != "jack" {
		t.Errorf("expected name to be 'jack', got '%s'", data2.Name)
	}

	err = c.Tag(tags...).Flush(context.Background())
	if err != nil {
		t.Error(err)
	}

	data3 := &Data{}
	err = c.Get(context.Background(), key, data3)
	if err != nil {
		t.Error(err)
	}
	if data3.Name == "jack" {
		t.Errorf("expected name to be '', got '%s'", data3.Name)
	}
}
