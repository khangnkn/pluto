package project

import (
	"github.com/nkhang/pluto/internal/rediskey"
	"github.com/nkhang/pluto/pkg/cache"
	"github.com/nkhang/pluto/pkg/errors"
	"github.com/nkhang/pluto/pkg/logger"
)

type Repository interface {
	Get(id uint64) (Project, error)
}

type repository struct {
	disk  DiskRepository
	cache cache.Cache
}

func NewRepository(r DiskRepository, c cache.Cache) *repository {
	return &repository{
		disk:  r,
		cache: c,
	}
}

func (r *repository) Get(id uint64) (Project, error) {
	var p Project
	k := rediskey.ProjectByID(id)
	err := r.cache.Get(k, &p)
	if err == nil {
		return p, nil
	}
	if errors.Type(err) != errors.CacheNotFound {
		logger.Error("error getting project from cache", err)
	} else {
		logger.Infof("cache miss for getting project: %+v", p)
	}
	p, err = r.disk.Get(id)
	if err != nil {
		return p, err
	}
	go func() {
		r.cache.Set(k, &p)
	}()
	return p, nil
}
