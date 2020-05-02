package tool

import (
	"github.com/nkhang/pluto/internal/rediskey"
	"github.com/nkhang/pluto/pkg/cache"
	"github.com/nkhang/pluto/pkg/errors"
	"github.com/nkhang/pluto/pkg/logger"
)

type Repository interface {
	GetAll() ([]Tool, error)
}

type repository struct {
	dbRepo    DbRepository
	cacheRepo cache.Cache
}

func NewRepository(d DbRepository, c cache.Cache) *repository {
	return &repository{
		dbRepo:    d,
		cacheRepo: c,
	}
}

func (r *repository) GetAll() ([]Tool, error) {
	tools := make([]Tool, 0)
	key := rediskey.AllTools()
	err := r.cacheRepo.Get(key, &tools)
	if err == nil {
		logger.Info("cache hit for all tools")
		return tools, nil
	}

	if errors.Type(err) == errors.CacheNotFound {
		logger.Infof("cache miss for all tools, key %s", key)
	} else {
		logger.Error("cannot get from cache")
	}

	tools, err = r.dbRepo.GetAll()
	if err != nil {
		logger.Error("cannot get from disk", err)
		return nil, err
	}
	go func() {
		err := r.cacheRepo.Set(rediskey.AllTools(), tools)
		if err != nil {
			logger.Error("cannot set all tools", err)
		}
	}()
	return tools, nil
}
