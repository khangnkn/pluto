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
	CreatePermission(workspaceID uint64, userIDs []uint64, role Role) error
	Create(userID uint64, title, description, color string) (Workspace, error)
	InvalidateWorkspacesForUser(userID uint64)
	InvalidatePermissionsForWorkspace(workspaceID uint64)
	UpdateWorkspace(workspaceID uint64, changes map[string]interface{}) (Workspace, error)
	DeleteWorkspace(workspaceID uint64) error
	DeletePermission(workspaceID uint64, userID uint64) error
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
	var (
		workspaces = make([]Workspace, 0)
		total      int
	)
	k, totalKey := rediskey.WorkspacesByUserID(userID, int32(role), offset, limit)
	err := r.cacheRepo.Get(k, &workspaces)
	err2 := r.cacheRepo.Get(totalKey, &total)
	if err == nil && err2 == nil {
		return workspaces, total, nil
	}
	if errors.Type(err) == errors.CacheNotFound {
		logger.Infof("cache miss for user %d", userID)
	} else {
		logger.Errorf("error getting cache workspaces for user %d", userID)
	}
	workspaces, total, err = r.dbRepo.GetByUserID(userID, role, offset, limit)
	if err != nil {
		logger.Error("error getting workspaces from database", err)
		return nil, 0, err
	}
	logger.Infof("getting workspace for user %d successfully", userID)
	return workspaces, total, nil
}

func (r *repository) GetPermission(workspaceID uint64, role Role, offset, limit int) ([]Permission, int, error) {
	var perms = make([]Permission, 0)
	var total int
	k, totalKey := rediskey.WorkspacesPermissionByWorkspaceID(workspaceID, int32(role), offset, limit)
	err := r.cacheRepo.Get(k, &perms)
	err2 := r.cacheRepo.Get(totalKey, &total)
	if err == nil && err2 == nil {
		logger.Infof("cache hit getting workspace permission %d", workspaceID)
		return perms, 0, nil
	}
	if errors.Type(err) == errors.CacheNotFound {
		logger.Infof("cache miss for workspace %d", workspaceID)
	} else {
		logger.Errorf("error getting cache perms for workspace %d", workspaceID)
	}
	perms, total, err = r.dbRepo.GetPermissionByWorkspaceID(workspaceID, role, offset, limit)
	if err != nil {
		logger.Error("error getting perms from database", err)
		return nil, 0, err
	}
	go func() {
		err := r.cacheRepo.Set(k, &perms)
		if err != nil {
			logger.Error(err)
		}
	}()
	logger.Infof("getting permission for workspace %d successfully", workspaceID)
	return perms, total, nil
}

func (r *repository) Create(userID uint64, title, description, color string) (Workspace, error) {

	w, err := r.dbRepo.Create(userID, title, description, color)
	if err != nil {
		return Workspace{}, err
	}
	go func() {
		r.InvalidateWorkspacesForUser(userID)
	}()
	return w, nil
}

func (r *repository) InvalidateWorkspacesForUser(userID uint64) {
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

func (r *repository) InvalidatePermissionsForWorkspace(workspaceID uint64) {
	pattern := rediskey.WorkspacesPermissionByWorkspaceIDPattern(workspaceID)
	keys, err := r.cacheRepo.Keys(pattern)
	if err != nil {
		logger.Errorf("cannot get all workspace permission by workspace id %v", err.Error())
		return
	}
	if err := r.cacheRepo.Del(keys...); err != nil {
		logger.Errorf("error delete all permission for workspace %d in redis", workspaceID)
	}
}

func (r *repository) UpdateWorkspace(workspaceID uint64, changes map[string]interface{}) (Workspace, error) {
	k := rediskey.WorkspaceByID(workspaceID)
	err := r.cacheRepo.Del(k)
	if err != nil {
		logger.Error(err)
	}
	return r.dbRepo.UpdateWorkspace(workspaceID, changes)
}

func (r *repository) CreatePermission(workspaceID uint64, userIDs []uint64, role Role) error {
	err := r.dbRepo.CreatePermission(workspaceID, userIDs, role)
	if err != nil {
		return err
	}
	go func() {
		r.InvalidatePermissionsForWorkspace(workspaceID)
		for i := range userIDs {
			r.InvalidateWorkspacesForUser(userIDs[i])
		}
	}()
	return nil
}

func (r *repository) DeletePermission(workspaceID uint64, userID uint64) error {
	err := r.dbRepo.DeletePermission(workspaceID, userID)
	if err != nil {
		return err
	}
	go func() {
		r.InvalidatePermissionsForWorkspace(workspaceID)
		r.InvalidateWorkspacesForUser(userID)
	}()
	return nil
}

func (r *repository) DeleteWorkspace(workspaceID uint64) error {
	perms, _, err := r.GetPermission(workspaceID, Any, 0, 0)
	if err != nil {
		return err
	}
	err = r.dbRepo.DeleteWorkspace(workspaceID)
	if err != nil {
		return err
	}
	go func() {
		r.InvalidatePermissionsForWorkspace(workspaceID)
		for i := range perms {
			r.InvalidateWorkspacesForUser(perms[i].UserID)
		}
	}()
	return nil
}
