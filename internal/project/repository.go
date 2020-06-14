package project

import (
	"github.com/nkhang/pluto/internal/rediskey"
	"github.com/nkhang/pluto/pkg/cache"
	"github.com/nkhang/pluto/pkg/errors"
	"github.com/nkhang/pluto/pkg/logger"
)

type Repository interface {
	Get(pID uint64) (Project, error)
	GetByWorkspaceID(id uint64) ([]Project, error)
	GetProjectPermission(pID uint64) ([]Permission, error)
	CreateProject(wID uint64, title, desc string) (Project, error)
	InvalidateProjectsByWorkspaceID(id uint64) error
}

type repository struct {
	disk  DBRepository
	cache cache.Cache
}

func NewRepository(r DBRepository, c cache.Cache) *repository {
	return &repository{
		disk:  r,
		cache: c,
	}
}

func (r *repository) Get(pID uint64) (Project, error) {
	var p Project
	k := rediskey.ProjectByID(pID)
	err := r.cache.Get(k, &p)
	if err == nil {
		logger.Infof("cache hit for getting project %d", pID)
		return p, nil
	}
	if errors.Type(err) != errors.CacheNotFound {
		logger.Error("error getting project from cache", err)
	} else {
		logger.Infof("cache miss for getting project %d", pID)
	}
	p, err = r.disk.Get(pID)
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
		logger.Infof("cache hit for getting projects for workspace %d", id)
		return projects, nil
	}
	if errors.Type(err) == errors.CacheNotFound {
		logger.Infof("cache miss for getting projects for workspace %d", id)
	} else {
		logger.Errorf("cannot get projects for workspace %d", id)
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

func (r *repository) CreateProject(wID uint64, title, desc string) (Project, error) {
	return r.disk.CreateProject(wID, title, desc)
}

func (r *repository) GetProjectPermission(pID uint64) ([]Permission, error) {
	var perms []Permission
	k := rediskey.ProjectPermissionByID(pID)
	err := r.cache.Get(k, &perms)
	if err == nil {
		logger.Infof("cache hit for getting permissions of projects %d", pID)
		return perms, nil
	}
	if errors.Type(err) == errors.CacheNotFound {
		logger.Infof("cache miss for getting permissions of projects %d", pID)
	} else {
		logger.Errorf("cannot get permissions of projects %d", pID)
	}
	perms, err = r.disk.GetProjectPermission(pID)
	if err != nil {
		return nil, err
	}
	go func() {
		err := r.cache.Set(k, &perms)
		if err != nil {
			logger.Error(err)
		}
	}()
	return perms, nil
}

func (r *repository) InvalidateProjectsByWorkspaceID(id uint64) error {
	k := rediskey.ProjectByWorkspaceID(id)
	return r.cache.Del(k)
}
