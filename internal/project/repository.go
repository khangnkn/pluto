package project

import (
	"github.com/nkhang/pluto/internal/dataset"
	"github.com/nkhang/pluto/internal/rediskey"
	"github.com/nkhang/pluto/internal/task"
	"github.com/nkhang/pluto/pkg/cache"
	"github.com/nkhang/pluto/pkg/errors"
	"github.com/nkhang/pluto/pkg/logger"
	uuid "github.com/satori/go.uuid"
)

type Repository interface {
	Get(pID uint64) (Project, error)
	CreateProject(wID uint64, title, desc, color string) (Project, error)
	GetByWorkspaceID(id uint64, offset, limit int) ([]Project, int, error)
	GetUserPermissions(userID uint64, role Role, offset, limit int) ([]Permission, int, error)
	GetProjectPermissions(pID uint64, role Role, offset, limit int) ([]Permission, int, error)
	GetPermission(userID, projectID uint64) (Permission, error)
	CreatePermission(projectID, userID uint64, role Role) (Permission, error)
	UpdatePermission(projectID, userID uint64, role Role) (Permission, error)
	UpdateProject(projectID uint64, changes map[string]interface{}) (Project, error)
	Delete(id uint64) error
	DeletePermission(userID, projectID uint64) error
	DeleteByWorkspace(workspaceID uint64) error
}

type repository struct {
	taskRepo    task.Repository
	datasetRepo dataset.Repository
	disk        DBRepository
	cache       cache.Cache
}

func NewRepository(r DBRepository, c cache.Cache, t task.Repository, d dataset.Repository) *repository {
	return &repository{
		taskRepo:    t,
		datasetRepo: d,
		disk:        r,
		cache:       c,
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
	r.invalidateProjectsByWorkspaceID(wID)
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
	r.invalidatePermissionForProject(projectID)
	r.invalidatePermissionForUser(userID)
	_, err := r.Get(projectID)
	if errors.Type(err) == errors.ProjectNotFound {
		return Permission{}, errors.ProjectNotFound.NewWithMessageF("project %d not existed", projectID)
	}
	return r.disk.CreatePermission(projectID, userID, role)
}

func (r *repository) UpdatePermission(projectID, userID uint64, role Role) (Permission, error) {
	perm, err := r.disk.UpdatePermission(projectID, userID, role)
	if err != nil {
		return Permission{}, err
	}
	r.invalidatePermissionForUser(userID)
	r.invalidatePermissionForProject(projectID)
	return perm, nil
}

func (r *repository) invalidateProjectsByWorkspaceID(id uint64) {
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

func (r *repository) invalidatePermissionForUser(userID uint64) {
	pattern := rediskey.ProjectPermissionByUserPattern(userID)
	keys, err := r.cache.Keys(pattern)
	if err != nil {
		logger.Errorf("error invalidate cache for user %d. err %v", userID, err.Error())
		return
	}
	err = r.cache.Del(keys...)
	if err != nil {
		logger.Errorf("error invalidate cache for user %d. err %v", userID, err.Error())
	}
}

func (r *repository) invalidatePermissionForProject(projectID uint64) {
	_, totalKey, pattern := rediskey.ProjectPermissionByID(projectID, 0, 0, 0)
	keys, err := r.cache.Keys(pattern)
	if err != nil {
		logger.Errorf("error getting all keys with pattern %s", pattern)
		return
	}
	projectKey := rediskey.ProjectByID(projectID)
	if err := r.cache.Del(append(keys, totalKey, projectKey)...); err != nil {
		logger.Errorf("error delete key %s", keys)
	}
	logger.Infof("invalidate key %d successfully", len(keys))
}

func (r *repository) GetPermission(userID, projectID uint64) (Permission, error) {
	return r.disk.GetPermission(userID, projectID)
}

func (r *repository) UpdateProject(projectID uint64, changes map[string]interface{}) (Project, error) {
	r.invalidateProject(projectID)
	project, err := r.disk.UpdateProject(projectID, changes)
	if err != nil {
		return project, errors.ProjectCannotUpdate.Wrap(err, "cannot update project")
	}
	go func() {
		r.invalidateProjectsByWorkspaceID(project.WorkspaceID)
		perms, _, err := r.GetProjectPermissions(projectID, Any, 0, 0)
		if err != nil {
			return
		}
		for _, v := range perms {
			r.invalidatePermissionForUser(v.UserID)
		}
	}()
	return project, nil
}

func (r *repository) Delete(id uint64) error {
	project, err := r.Get(id)
	if err != nil {
		return err
	}
	perms, _, err := r.GetProjectPermissions(id, Any, 0, 0)
	if err != nil {
		return err
	}
	r.invalidateProjectsByWorkspaceID(project.WorkspaceID)
	r.invalidatePermissionForProject(id)
	r.invalidateProject(id)
	err = r.disk.Delete(id)

	for _, v := range perms {
		r.invalidatePermissionForUser(v.UserID)
	}
	err = r.taskRepo.DeleteTaskByProject(id)
	if err != nil {
		logger.Errorf("error deleting tasks of project %d. err %v", id, err)
	}
	err = r.datasetRepo.DeleteByProject(id)
	if err != nil {
		logger.Errorf("error deleting datasets of project %d, err %v", id, err)
	}
	return nil
}

func (r *repository) DeleteByWorkspace(workspaceID uint64) error {
	r.invalidateProjectsByWorkspaceID(workspaceID)
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
func (r *repository) DeletePermission(userID, projectID uint64) error {
	err := r.disk.DeletePermission(userID, projectID)
	if err != nil {
		return err
	}
	r.invalidatePermissionForUser(userID)
	r.invalidatePermissionForProject(projectID)
	return nil
}

func (r *repository) PickThumbnail() {

}

func (r *repository) invalidateProject(projectID uint64) {
	k := rediskey.ProjectByID(projectID)
	err := r.cache.Del(k)
	if err != nil {
		logger.Error(err)
	}
}
