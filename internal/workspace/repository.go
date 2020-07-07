package workspace

import (
	"github.com/nkhang/pluto/internal/rediskey"
	"github.com/nkhang/pluto/pkg/cache"
	"github.com/nkhang/pluto/pkg/errors"
	"github.com/nkhang/pluto/pkg/logger"
)

type Repository interface {
	Get(id uint64) (Workspace, error)
	GetByUserID(userID uint64, role Role, offset, limit int) ([]Workspace, int, error)
	GetPermission(workspaceID uint64, role Role, offset, limit int) ([]Permission, int, error)
	Create(userID uint64, title, description string) (Workspace, error)
	InvalidateForUser(userID uint64)
}

type repository struct {
	dbRepo    DBRepository
	cacheRepo cache.Cache
}

func NewRepository(dbRepo DBRepository, c cache.Cache) *repository {
	return &repository{
		dbRepo:    dbRepo,
		cacheRepo: c,
	}
}

func (r *repository) Get(id uint64) (Workspace, error) {
	var w Workspace
	k := rediskey.WorkspaceByID(id)
	err := r.cacheRepo.Get(k, &w)
	if err == nil {
		return w, nil
	}
	if errors.Type(err) == errors.CacheNotFound {
		logger.Infof("cache miss for project %d", id)
	} else {
		logger.Errorf("error getting cache for workspace %d", id)
	}
	w, err = r.dbRepo.Get(id)
	if err != nil {
		return Workspace{}, err
	}
	logger.Infof("getting workspace [%d] successfully", id)
	go func() {
		err := r.cacheRepo.Set(k, &w)
		if err != nil {
			logger.Error("error setting cache of workspace")
		}
	}()
	return w, nil
}

func (r *repository) GetByUserID(userID uint64, role Role, offset, limit int) ([]Workspace, int, error) {
	var workspaces = make([]Workspace, 0)
	k := rediskey.WorkspacesByUserID(userID, int32(role), offset, limit)
	err := r.cacheRepo.Get(k, &workspaces)
	if err == nil {
		return workspaces, 0, nil
	}
	if errors.Type(err) == errors.CacheNotFound {
		logger.Infof("cache miss for user %d", userID)
	} else {
		logger.Errorf("error getting cache workspaces for user %d", userID)
	}
	workspaces, total, err := r.dbRepo.GetByUserID(userID, role, offset, limit)
	if err != nil {
		logger.Error("error getting workspaces from database", err)
		return nil, 0, err
	}
	logger.Infof("getting workspace for user %d successfully", userID)
	return workspaces, total, nil
}

func (r *repository) GetPermission(workspaceID uint64, role Role, offset, limit int) ([]Permission, int, error) {
	var perms = make([]Permission, 0)
	k := rediskey.WorkspacesPermissionByWorkspaceID(workspaceID, int32(role), offset, limit)
	err := r.cacheRepo.Get(k, &perms)
	if err == nil {
		return perms, 0, nil
	}
	if errors.Type(err) == errors.CacheNotFound {
		logger.Infof("cache miss for workspace %d", workspaceID)
	} else {
		logger.Errorf("error getting cache perms for workspace %d", workspaceID)
	}
	perms, total, err := r.dbRepo.GetPermissionByWorkspaceID(workspaceID, role, offset, limit)
	if err != nil {
		logger.Error("error getting perms from database", err)
		return nil, 0, err
	}
	logger.Infof("getting permission for workspace %d successfully", workspaceID)
	return perms, total, nil
}

func (r *repository) Create(userID uint64, title, description string) (Workspace, error) {
	return r.dbRepo.Create(userID, title, description)
}

func (r *repository) InvalidateForUser(userID uint64) {
	pattern := rediskey.WorkspacesByUserIDPattern(userID)
	keys, err := r.cacheRepo.Keys(pattern)
	if err != nil {
		logger.Error("error getting keys to invalidate for user")
		return
	}
	if len(keys) == 0 {
		return
	}
	err = r.cacheRepo.Del(keys...)
	if err != nil {
		logger.Errorf("error delete keys %v", keys)
		return
	}
}
