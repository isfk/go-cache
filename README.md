# cache
go cache based on go-redis

## Installation

Install:
```bash
go get -u github.com/go-cache/cache
```
Import:
```bash
import "github.com/go-cache/cache"
```

## QuickStart

### Init with `cache.NewClient(conf)`
```gotemplate
package model

import (
	"github.com/go-cache/cache"
	"github.com/go-redis/redis/v7"
)

// CacheDriver CacheDriver
var CacheDriver *cache.Client

// InitCache InitCache
func InitCache() {
	conf := redis.Options{
		Addr: "redis:6379",
		Password: "sdfsdf",
	}
	CacheDriver = cache.NewClient(conf)
}
```

```gotemplate
InitCache()
```

### Command
All redis command check: https://godoc.org/github.com/go-redis/redis
```gotemplate
CacheDriver.RedisClient.Set("key", "value", time.Hour).Err()
CacheDriver.RedisClient.Get("key").Result()
CacheDriver.RedisClient.command()...
```

### use `Tag` & `Put`
`Put` is the same as `Set`, but just with `Tag` together.
```gotemplate
CacheDriver.Tag("user_all", "user_list").Put("user_id:1", &proto.User{Id: 1, Nickname: "111"}, time.Hour)
CacheDriver.Tag("user_all").Put("user_id:2", &proto.User{Id: 2, Nickname: "222"}, time.Hour)
CacheDriver.Tag("user_list").Clear()
```
