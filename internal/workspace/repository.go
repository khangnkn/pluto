package workspace

import (
	"github.com/nkhang/pluto/internal/project"
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
	UpdateWorkspace(workspaceID uint64, changes map[string]interface{}) (Workspace, error)
	DeleteWorkspace(workspaceID uint64) error
	DeletePermission(workspaceID uint64, userID uint64) error
}

type repository struct {
	dbRepo      DBRepository
	projectRepo project.Repository
	cacheRepo   cache.Cache
}

func NewRepository(dbRepo DBRepository, projectRepo project.Repository, c cache.Cache) *repository {
	return &repository{
		dbRepo:      dbRepo,
		projectRepo: projectRepo,
		cacheRepo:   c,
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
	k, totalKey, _ := rediskey.WorkspacesByUserID(userID, int32(role), offset, limit)
	err := r.cacheRepo.Get(k, &workspaces)
	err2 := r.cacheRepo.Get(totalKey, &total)
	if err == nil && err2 == nil {
		logger.Infof("cache hit for getting workspace by user ID %d", userID)
		return workspaces, total, nil
	}
	if errors.Type(err) == errors.CacheNotFound {
		logger.Infof("cache miss for user %d", userID)
	} else {
		logger.Infof("error getting cache workspaces for user %d. error %v", userID, err)
	}
	workspaces, total, err = r.dbRepo.GetByUserID(userID, role, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	go func() {
		err1 := r.cacheRepo.Set(k, workspaces)
		err2 := r.cacheRepo.Set(totalKey, total)
		if err1 != nil || err2 != nil {
			logger.Error("error setting workspace to cache")
		}
	}()
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
		return perms, total, nil
	}
	if errors.Type(err) == errors.CacheNotFound {
		logger.Infof("cache miss for workspace %d", workspaceID)
	} else {
		logger.Infof("error getting cache perms for workspace %d, error %v", workspaceID, err)
	}
	perms, total, err = r.dbRepo.GetPermissionByWorkspaceID(workspaceID, role, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	go func() {
		err := r.cacheRepo.Set(k, &perms)
		if err != nil {
			logger.Error(err)
		}
		logger.Infof("total keys %d", total)
		err = r.cacheRepo.Set(totalKey, total)
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
	err = r.CreatePermission(w.ID, []uint64{userID}, Admin)
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
	logger.Infof("invalidate workspaces permission for users %d successfully", userID)
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
	_, err := r.dbRepo.Get(workspaceID)
	if errors.Type(err) == errors.WorkspaceNotFound {
		return err
	}
	for _, userID := range userIDs {
		_, err = r.dbRepo.GetPermission(workspaceID, userID)
		if err == nil {
			continue
		}
		err = r.dbRepo.CreatePermission(workspaceID, userID, role)
	}
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
	err := r.triggerDeleteProjectsPermissions(userID, workspaceID)
	if err != nil {
		return err
	}
	err = r.dbRepo.DeletePermission(workspaceID, userID)
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
	err = r.projectRepo.DeleteByWorkspace(workspaceID)
	if err != nil {
		logger.Errorf("[WORKSPACE] error delete projects by workspace %d. err %v", workspaceID, err)
	}
	go func() {
		r.InvalidatePermissionsForWorkspace(workspaceID)
		for i := range perms {
			r.InvalidateWorkspacesForUser(perms[i].UserID)
		}
	}()
	return nil
}

func (r *repository) triggerDeleteProjectsPermissions(userID, workspaceID uint64) (err error) {
	projects, err := r.projectRepo.GetByWorkspaceID(workspaceID)
	if err != nil {
		return
	}
	for i := range projects {
		err := r.projectRepo.DeletePermission(userID, projects[i].ID)
		if err != nil {
			logger.Errorf("[WORKSPACE] - error deleting project permission for user %d in project %d", userID, projects[i].ID)
		}
	}
	return nil
}
