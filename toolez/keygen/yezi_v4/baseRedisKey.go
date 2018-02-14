package yezi_v4

import (
	"time"

	"github.com/go-redis/redis"
	"github.com/vettu/allez/libs/errorf"
)

type baseRedisKey struct {
	Conn     *redis.Client
	Key      string
	Duration time.Duration

	rawValue string
}

func (brk *baseRedisKey) checkKey() (err error) {
	if brk.Key == "" {
		err = errorf.New("redis key not set")
	}
	return
}

func (brk *baseRedisKey) Get() (err error) {
	if err = brk.checkKey(); err != nil {
		return
	}
	brk.rawValue, err = brk.Conn.Get(brk.Key).Result()
	return
}

func (brk *baseRedisKey) Set() (err error) {
	if err = brk.checkKey(); err != nil {
		return
	}
	err = brk.Conn.Set(brk.Key, brk.rawValue, brk.Duration).Err()
	return
}

func (brk *baseRedisKey) GetKey() string {
	return brk.Key
}
