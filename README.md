# cache
cache based on go-redis

## Installation

Install:
```bash
go get -u github.com/go-cache/cache/v2
```
Import:
```bash
import "github.com/go-cache/cache/v2"
```

## QuickStart

### Init with `cache.NewClient(conf)`

```gotemplate
rdb := redisV8.NewClient(&redisV8.Options{
    Addr:     redis.StdRedisConfig("main").Addr,
    Password: redis.StdRedisConfig("main").Password, // no password set
    DB:       redis.StdRedisConfig("main").DB,
})

CacheModel = cache.NewClient(context.Background(), rdb)
```

### Command

- use `Tag` & `Set` & `Get`

`Set` is the same as `Set`, but just with `Tag` together.

```gotemplate
CacheModel.Tag("user_all", "user_list").Set("user_id:1", &proto.User{Id: 1, Nickname: "111"}, time.Hour)
CacheModel.Tag("user_all").Set("user_id:2", &proto.User{Id: 2, Nickname: "222"}, time.Hour)
CacheModel.Get("user_id:2")
CacheModel.Tag("user_list").Flush()
```

- use redis

All redis command check: https://godoc.org/github.com/go-redis/redis

```gotemplate
CacheModel.RedisClient.Set("key", "value", time.Hour).Err()
CacheModel.RedisClient.Get("key").Result()
CacheModel.RedisClient.command()...
```

