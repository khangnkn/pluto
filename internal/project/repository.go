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
	GetByWorkspaceID(id uint64, offset, limit int) ([]Project, int, error)
	GetUserPermissions(userID uint64, role Role, offset, limit int) ([]Permission, int, error)
	GetProjectPermissions(pID uint64, role Role, offset, limit int) ([]Permission, int, error)
	GetPermission(userID, projectID uint64) (Permission, error)
	CreateProject(wID uint64, title, desc, color string) (Project, error)
	CreatePermission(projectID, userID uint64, role Role) (Permission, error)
	UpdateProject(projectID uint64, changes map[string]interface{}) (Project, error)
	Delete(id uint64) error
	DeleteByWorkspace(workspaceID uint64) error
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
		logger.Errorf("error getting project from cache %v", err)
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

func (r *repository) GetByWorkspaceID(id uint64, offset, limit int) ([]Project, int, error) {
	var projects = make([]Project, 0)
	var total int
	k, totalKey, _ := rediskey.ProjectByWorkspaceID(id, offset, limit)
	err := r.cache.Get(k, &projects)
	err2 := r.cache.Get(totalKey, &total)
	if err == nil && err2 == nil {
		logger.Infof("cache hit for getting projects for workspace %d", id)
		return projects, total, nil
	}
	if errors.Type(err) == errors.CacheNotFound {
		logger.Infof("cache miss for getting projects for workspace %d", id)
	} else {
		logger.Errorf("cannot get projects for workspace %d", id)
	}
	projects, total, err = r.disk.GetByWorkspaceID(id, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	go func() {
		err := r.cache.Set(k, &projects)
		if err != nil {
			logger.Error(err)
		}
		err = r.cache.Set(totalKey, &total)
		if err != nil {
			logger.Error(err)
		}
	}()
	return projects, total, nil
}

func (r *repository) GetUserPermissions(userID uint64, role Role, offset, limit int) (p []Permission, total int, err error) {
	k, totalKey := rediskey.PermissionsByUserID(userID, int32(role), offset, limit)
	err = r.cache.Get(k, &p)
	err2 := r.cache.Get(totalKey, &total)
	if err == nil && err2 == nil {
		logger.Infof("cache hit for getting user project permission, total %d perms", len(p))
		return
	}
	if errors.Type(err) == errors.CacheNotFound {
		logger.Info("cache miss for getting user projects")
	} else {
		logger.Infof("error getting user projects. error %s", err.Error())
	}
	p, total, err = r.disk.GetUserPermissions(userID, role, offset, limit)
	if err != nil {
		return
	}
	go func() {
		err := r.cache.Set(k, p)
		if err != nil {
			logger.Infof("error setting cache for get user projects")
		}
		err = r.cache.Set(totalKey, total)
		if err != nil {
			logger.Infof("error setting cache for get user projects")
		}
	}()
	return
}

func (r *repository) CreateProject(wID uint64, title, desc, color string) (Project, error) {
	r.InvalidateProjectsByWorkspaceID(wID)
	uid := uuid.NewV4().String()
	return r.disk.CreateProject(wID, title, desc, color, uid)
}

func (r *repository) GetProjectPermissions(pID uint64, role Role, offset, limit int) ([]Permission, int, error) {
	var (
		perms []Permission
		total int
	)
	specKey, totalKey, _ := rediskey.ProjectPermissionByID(pID, uint32(role), offset, limit)
	err := r.cache.Get(specKey, &perms)
	err2 := r.cache.Get(totalKey, &total)
	if err == nil && err2 == nil {
		logger.Infof("cache hit for getting permissions of projects %d", pID)
		return perms, total, nil
	}
	if errors.Type(err) == errors.CacheNotFound {
		logger.Infof("cache miss for getting permissions of projects %d", pID)
	} else {
		logger.Errorf("cannot get permissions of projects %d", pID)
	}
	perms, total, err = r.disk.GetProjectPermissions(pID, role, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	go func() {
		err := r.cache.Set(specKey, &perms)
		if err != nil {
			logger.Error(err)
		}
		err = r.cache.Set(totalKey, total)
		if err != nil {
			logger.Error(err)
		}
	}()
	return perms, total, nil
}

func (r *repository) CreatePermission(projectID, userID uint64, role Role) (Permission, error) {
	r.InvalidatePermissionForProject(projectID)
	r.InvalidatePermissionForUser(userID)
	_, err := r.Get(projectID)
	if errors.Type(err) == errors.ProjectNotFound {
		return Permission{}, errors.ProjectNotFound.NewWithMessageF("project %d not existed", projectID)
	}
	return r.disk.CreatePermission(projectID, userID, role)
}

func (r *repository) InvalidateProjectsByWorkspaceID(id uint64) {
	_, totalKey, pattern := rediskey.ProjectByWorkspaceID(id, 0, 0)
	keys, err := r.cache.Keys(pattern)
	if err != nil {
		logger.Errorf("error getting pattern %s", pattern)
	}
	keys = append(keys, totalKey)
	if err := r.cache.Del(keys...); err != nil {
		logger.Errorf("error deleting keys %v", keys)
	}
}

func (r *repository) InvalidatePermissionForUser(userID uint64) error {
	_, totalKey, pattern := rediskey.ProjectByWorkspaceID(userID, 0, 0)
	keys, err := r.cache.Keys(pattern)
	if err != nil {
		return err
	}
	keys = append(keys, totalKey)
	return r.cache.Del(keys...)
}

func (r *repository) InvalidatePermissionForProject(projectID uint64) {
	_, _, pattern := rediskey.ProjectPermissionByID(projectID, 0, 0, 0)
	keys, err := r.cache.Keys(pattern)
	if err != nil {
		logger.Errorf("error getting all keys with pattern %s", pattern)
		return
	}
	if err := r.cache.Del(keys...); err != nil {
		logger.Errorf("error delete key %s", keys)
	}
	logger.Infof("invalidate key %d successfully", len(keys))
}

func (r *repository) GetPermission(userID, projectID uint64) (Permission, error) {
	return r.disk.GetPermission(userID, projectID)
}

func (r *repository) UpdateProject(projectID uint64, changes map[string]interface{}) (Project, error) {
	k := rediskey.ProjectByID(projectID)
	err := r.cache.Del(k)
	if err != nil {
		logger.Error(err)
	}
	project, err := r.disk.UpdateProject(projectID, changes)
	if err != nil {
		return project, errors.ProjectCannotUpdate.Wrap(err, "cannot update project")
	}
	go func() {
		r.InvalidateProjectsByWorkspaceID(project.WorkspaceID)
	}()
	return project, nil
}

func (r *repository) Delete(id uint64) error {
	project, err := r.Get(id)
	if err != nil {
		return err
	}
	r.InvalidateProjectsByWorkspaceID(project.WorkspaceID)
	r.InvalidatePermissionForProject(id)
	return r.disk.Delete(id)
}

func (r *repository) DeleteByWorkspace(workspaceID uint64) error {
	r.InvalidateProjectsByWorkspaceID(workspaceID)
	projects, _, err := r.GetByWorkspaceID(workspaceID, 0, 0)
	if err != nil {
		return err
	}
	for i := range projects {
		err = r.Delete(projects[i].ID)
		if err != nil {
			logger.Errorf("error delete project %d of workspace %d", projects[i].ID, workspaceID)
		}
	}
	return nil
}
