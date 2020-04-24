package tool

import (
	"github.com/nkhang/pluto/internal/rediskey"
	"github.com/nkhang/pluto/pkg/cache"
	"github.com/nkhang/pluto/pkg/errors"
	"github.com/nkhang/pluto/pkg/logger"
)

type repository struct {
	disk  *diskRepository
	cache cache.Cache
}

func NewRepository(d *diskRepository, c cache.Cache) *repository {
	return &repository{
		disk:  d,
		cache: c,
	}
}

func (r *repository) GetAll() ([]Tool, error) {
	tools := make([]Tool, 0)
	key := rediskey.AllTools()
	err := r.cache.Get(key, &tools)
	if errors.Type(err) == errors.CacheNotFound {
		logger.Infof("cache miss for all tools, key %s", key)
		err = nil
	}
	if err == nil {
		logger.Info("cache hit for all tools")
	} else {
		logger.Error("cannot get from cache")
	}
	tools, err = r.disk.GetAll()
	if err != nil {
		logger.Error("cannot get from disk", err)
		return nil, err
	}
	go func() {
		err := r.cache.Set(rediskey.AllTools(), tools)
		if err != nil {
			logger.Error("cannot set all tools", err)
		}
	}()
	return tools, nil
}
