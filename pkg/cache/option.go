package cache

import "time"

type Option interface {
	apply(c *client)
}

type optionFn func(c *client)

func (f optionFn) apply(c *client) {
	f(c)
}

func WithExpireTime(t time.Duration) optionFn {
	return func(c *client) {
		c.exp = t
	}
}
