# cache

cache based on go-redis

## Installation

Install:

```bash
go get -u github.com/go-cache/cache/v3
```

Import:

```bash
import "github.com/go-cache/cache/v3"
```

## QuickStart

### Init with `cache.New(conf)`

```gotemplate
rdb := redisV9.NewClient(&redisV9.Options{
    Addr:     redis.StdRedisConfig("main").Addr,
    Password: redis.StdRedisConfig("main").Password, // no password set
    DB:       redis.StdRedisConfig("main").DB,
})

c = cache.New(context.Background(), rdb, cache.WithPrefix("test"), cache.WithExpired(600*time.Second))
```

### Command

- use `Tag` & `Set` & `Get`

`Set` is the same as `Set`, but just with `Tag` together.

`Get` without `Tag`.

```gotemplate
c.Tag("tag:user:all", "tag:user:1").Set(ctx, "key:user:1", &proto.User{Id: 1, Nickname: "111"})
c.Tag("tag:user:all", "tag:user:2").Set(ctx, "key:user:2", &proto.User{Id: 2, Nickname: "222"})
c.Get(ctx, "key:user:1", &proto.User{})
c.Get(ctx, "key:user:2", &proto.User{})
c.Tag("tag:user:all").Flush()
```

- use redis

All redis command check: https://godoc.org/github.com/go-redis/redis

```gotemplate
c.redis.Set("key", "value", time.Hour).Err()
c.redis.Get("key").Result()
c.redis.command()...
```
