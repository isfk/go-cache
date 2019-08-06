package cache

import (
	"github.com/go-redis/cache"
	"github.com/go-redis/redis"
	"github.com/vmihailenco/msgpack"
)

// Driver Driver
var Driver *cache.Codec

// GetCodec GetCodec
func GetCodec() {
	ring := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"server1": "redis:6379",
		},
		Password: "sdfsdf",
	})

	codec := &cache.Codec{
		Redis: ring,

		Marshal: func(v interface{}) ([]byte, error) {
			return msgpack.Marshal(v)
		},
		Unmarshal: func(b []byte, v interface{}) error {
			return msgpack.Unmarshal(b, v)
		},
	}
	Driver = codec
}
