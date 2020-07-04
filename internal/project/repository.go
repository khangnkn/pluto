package project

import (
	"github.com/nkhang/pluto/internal/rediskey"
	"github.com/nkhang/pluto/pkg/cache"
	"github.com/nkhang/pluto/pkg/errors"
	"github.com/nkhang/pluto/pkg/logger"
	uuid "github.com/satori/go.uuid"
)

type Repository interface {
	Get(pID uint64) (Project, error)
	GetByWorkspaceID(id uint64) ([]Project, error)
	GetUserPermissions(userID uint64, role Role, offset, limit int) ([]Permission, error)
	GetProjectPermissions(pID uint64) ([]Permission, error)
	GetPermission(userID, projectID uint64) (Permission, error)
	CreateProject(wID uint64, title, desc string) (Project, error)
	CreatePermission(projectID, userID uint64, role Role) (Permission, error)
	InvalidateProjectsByWorkspaceID(id uint64) error
	InvalidatePermissionForUser(userID uint64) error
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

func (r *repository) GetUserPermissions(userID uint64, role Role, offset, limit int) (p []Permission, err error) {
	k := rediskey.PermissionsByUserID(userID, int32(role), offset, limit)
	err = r.cache.Get(k, &p)
	if err == nil {
		logger.Infof("cache hit for getting user project permission, total %d perms", len(p))
		return
	}
	if errors.Type(err) == errors.CacheNotFound {
		logger.Info("cache miss for getting user projects")
	} else {
		logger.Infof("error getting user projects. error %s", err.Error())
	}
	p, err = r.disk.GetUserPermissions(userID, role, offset, limit)
	if err != nil {
		return
	}
	go func() {
		err := r.cache.Set(k, p)
		if err != nil {
			logger.Infof("error setting cache for get user projects")
		}
	}()
	return
}

func (r *repository) CreateProject(wID uint64, title, desc string) (Project, error) {
	uid := uuid.NewV4().String()
	return r.disk.CreateProject(wID, title, desc, uid)
}

func (r *repository) GetProjectPermissions(pID uint64) ([]Permission, error) {
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
	perms, err = r.disk.GetProjectPermissions(pID)
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

func (r *repository) CreatePermission(projectID, userID uint64, role Role) (Permission, error) {
	return r.disk.CreatePermission(projectID, userID, role)
}

func (r *repository) InvalidateProjectsByWorkspaceID(id uint64) error {
	k := rediskey.ProjectByWorkspaceID(id)
	return r.cache.Del(k)
}

func (r *repository) InvalidatePermissionForUser(userID uint64) error {
	pattern := rediskey.ProjectPermissionByUserPattern(userID)
	keys, err := r.cache.Keys(pattern)
	if err != nil {
		return err
	}
	return r.cache.Del(keys...)
}

func (r *repository) GetPermission(userID, projectID uint64) (Permission, error) {
	return r.disk.GetPermission(userID, projectID)
}
