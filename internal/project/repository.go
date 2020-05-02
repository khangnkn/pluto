package project

import (
	"github.com/nkhang/pluto/internal/rediskey"
	"github.com/nkhang/pluto/pkg/cache"
	"github.com/nkhang/pluto/pkg/errors"
	"github.com/nkhang/pluto/pkg/logger"
)

type Repository interface {
	Get(wID uint64, pID uint64) (Project, error)
	GetByWorkspaceID(id uint64) ([]Project, error)
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

func (r *repository) Get(wID uint64, pID uint64) (Project, error) {
	var p Project
	k := rediskey.ProjectByID(wID, pID)
	err := r.cache.Get(k, &p)
	if err == nil {
		return p, nil
	}
	if errors.Type(err) != errors.CacheNotFound {
		logger.Error("error getting project from cache", err)
	} else {
		logger.Infof("cache miss for getting project: %+v", p)
	}
	p, err = r.disk.Get(wID, pID)
	if err != nil {
		return p, err
	}
	go func() {
		err := r.cache.Set(k, &p)
		if err != nil {
			logger.Error("error in setting cache", err)
		}
	}()
	return p, nil
}

func (r *repository) GetByWorkspaceID(id uint64) ([]Project, error) {
	var projects = make([]Project, 0)
	k := rediskey.ProjectByWorkspaceID(id)
	err := r.cache.Get(k, &projects)
	if err == nil {
		return projects, nil
	}
	if errors.Type(err) == errors.CacheNotFound {
		logger.Infof("cache miss for getting projects for workspace [%d]", id)
	} else {
		logger.Errorf("cannot get projects for workspace [%d]", id)
	}
	projects, err = r.disk.GetByWorkspaceID(id)
	if err != nil {
		return nil, err
	}
	go func() {
		err := r.cache.Set(k, &projects)
		if err != nil {
			logger.Error(err)
		}
	}()
	return projects, nil
}
