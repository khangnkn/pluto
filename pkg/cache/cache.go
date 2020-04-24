package cache

import (
	"time"

	"github.com/nkhang/pluto/pkg/errors"

	"github.com/go-redis/redis"
	"github.com/nkhang/pluto/pkg/json"
)

type Cache interface {
	Get(key string, target interface{}) error
	Set(key string, target interface{}) error
}

type client struct {
	cmd redis.UniversalClient
	exp time.Duration
}

func New(cmd redis.UniversalClient, opt ...Option) *client {
	c := &client{
		cmd: cmd,
		exp: 0,
	}
	for _, o := range opt {
		o.apply(c)
	}
	return c
}

func (c *client) Get(key string, target interface{}) error {
	data, err := c.cmd.Get(key).Result()
	if err == redis.Nil {
		return errors.CacheNotFound.NewWithMessage("cache not found")
	}
	if err != nil {
		return errors.CacheGetError.Wrap(err, "cannot get from cache")
	}
	return json.Deserialize(data, target)
}

func (c *client) Set(key string, data interface{}) error {
	data, err := json.Serialize(data)
	if err != nil {
		return errors.CacheSetError.Wrap(err, "cannot deserialize data")
	}
	err = c.cmd.Set(key, data, c.exp).Err()
	if err != nil {
		return errors.CacheSetError.Wrap(err, "cannot set cache")
	}
	return nil
}
