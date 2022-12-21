# cache

cache based on [go-redis](https://github.com/go-redis/redis)

## Installation

Install:

```bash
go get -u github.com/isfk/go-cache/v3
```

Import:

```bash
import "github.com/isfk/go-cache/v3"
```

## QuickStart

### Init with `cache.New()`

```gotemplate
rdb := redisV9.NewClient(&redisV9.Options{
    Addr:     "127.0.0.1:6379",
})

c := cache.New(
    context.Background(),
    rdb,
    cache.WithPrefix("test"),
    cache.WithExpired(600*time.Second),
)
```

### Command

- use `Tag` & `Set` & `Get`

`Set` is the same as `Set`, but just with `Tag` together.

`Get` without `Tag`.

```gotemplate
c.Tag([]string{"tag:user:all", "tag:user:1"}...).Set(ctx, "key:user:1", &proto.User{Id: 1, Nickname: "111"})
c.Tag([]string{"tag:user:all", "tag:user:2"}...).Set(ctx, "key:user:2", &proto.User{Id: 2, Nickname: "222"})

info := &proto.User{}
c.Get(ctx, "key:user:1", info)
c.Get(ctx, "key:user:2", info)
c.Tag([]string{"tag:user:all"}...).Flush(ctx)
```

- use redis

All redis command check: https://godoc.org/github.com/go-redis/redis/v9

```gotemplate
c.redis.Set("key", "value", time.Hour).Err()
c.redis.Get("key").Result()
c.redis.command()...
```
