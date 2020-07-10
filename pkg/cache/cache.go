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
	Del(key ...string) error
	Keys(pattern string) ([]string, error)
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

func (c *client) Del(key ...string) error {
	if len(key) == 0 {
		return nil
	}
	err := c.cmd.Del(key...).Err()
	if err != nil {
		return errors.CacheDeleteError.Wrap(err, "cannot del keys")
	}
	return nil
}

func (c *client) Keys(pattern string) ([]string, error) {
	s, err := c.cmd.Keys(pattern).Result()
	if err != nil {
		return nil, errors.CacheKeysError.Wrap(err, "cannot get key pattern")
	}
	return s, nil
}
