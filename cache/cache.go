package cache

import (
	"encoding/json"
	"errors"
	"time"

	cache "gopkg.in/go-redis/cache.v5"
	redis "gopkg.in/redis.v5"
	//msgpack "gopkg.in/vmihailenco/msgpack.v2"
)

var (
	ErrEmptyKey = errors.New("key is empty")
)

var codec *cache.Codec

//Configure set the cache config
func Configure(addr string, psw string, db int, poolSize int) {
	codec = &cache.Codec{
		Redis: redis.NewClient(&redis.Options{
			Network:  "tcp",
			Addr:     addr,
			Password: psw,
			DB:       db,
			PoolSize: poolSize,
		}),

		Marshal: func(v interface{}) ([]byte, error) {
			//return msgpack.Marshal(v)
			return json.Marshal(v)
		},

		Unmarshal: func(data []byte, v interface{}) error {
			//return msgpack.Unmarshal(data, v)
			return json.Unmarshal(data, v)
		},
	}
}

func Set(key string, v interface{}, expire time.Duration) error {
	if key == "" {
		return ErrEmptyKey
	}

	if err := codec.Set(&cache.Item{
		Key:        key,
		Object:     v,
		Expiration: expire,
	}); err != nil {
		return err
	}

	return nil
}

func Get(key string, v interface{}) error {
	if err := codec.Get(key, v); err != nil {
		return err
	}

	return nil
}

func Del(key string) error {
	if err := codec.Delete(key); err != nil {
		return err
	}

	return nil
}
