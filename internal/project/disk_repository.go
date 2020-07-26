package project

import (
	"github.com/jinzhu/gorm"
	"github.com/nkhang/pluto/pkg/logger"

	"github.com/nkhang/pluto/pkg/errors"
)

type DBRepository interface {
	Get(pID uint64) (Project, error)
	GetByWorkspaceID(wID uint64, offset, limit int) ([]Project, int, error)
	GetProjectPermissions(pID uint64, role Role, offset, limit int) (perms []Permission, total int, err error)
	GetUserPermissions(userID uint64, role Role, offset, limit int) ([]Permission, int, error)
	GetPermission(userID, projectID uint64) (Permission, error)
	CreateProject(wID uint64, title, desc, color, uid string) (Project, error)
	CreatePermission(projectID, userID uint64, role Role) (Permission, error)
	UpdatePermission(projectID, userID uint64, role Role) (Permission, error)
	UpdateProject(ProjectID uint64, changes map[string]interface{}) (Project, error)
	Delete(id uint64) error
	DeletePermission(userID, projectID uint64) error
}

type dbRepository struct {
	db *gorm.DB
}

func NewDiskRepository(db *gorm.DB) *dbRepository {
	return &dbRepository{db: db}
}

func (r *dbRepository) Get(pID uint64) (Project, error) {
	var p Project
	result := r.db.First(&p, pID)
	if result.RecordNotFound() {
		return Project{}, errors.ProjectNotFound.NewWithMessage("project not found")
	}
	if err := result.Error; err != nil {
		return Project{}, errors.ProjectQueryError.Wrap(err, "query project error")
	}
	return p, nil
}

func (r *dbRepository) GetByWorkspaceID(wID uint64, offset, limit int) ([]Project, int, error) {
	var projects = make([]Project, 0)
	var total int
	db := r.db.Model(&Project{}).
		Where(fieldWorkspaceID+" = ?", wID).
		Count(&total)
	if offset != 0 || limit != 0 {
		db = db.Offset(offset).Limit(limit)
	}
	err := db.Find(&projects).Error
	if err != nil {
		return nil, 0, errors.ProjectQueryError.NewWithMessage("error getting project of workspace")
	}
	return projects, total, nil
}

func (r *dbRepository) CreateProject(wID uint64, title, desc, color, uid string) (Project, error) {
	var p = Project{
		WorkspaceID: wID,
		Title:       title,
		Description: desc,
		Dir:         uid,
		Thumbnail:   defaultImage,
		Color:       color,
	}
	err := r.db.Create(&p).Error
	if err != nil {
		return Project{}, errors.ProjectCreatingError.Wrap(err, "cannot create project")
	}
	return p, nil
}

func (r *dbRepository) GetProjectPermissions(pID uint64, role Role, offset, limit int) (perms []Permission, total int, err error) {
	db := r.db.Model(&Permission{}).Where("project_id = ?", pID)
	if role != Any {
		db = db.Where("role = ?", role)
	}
	db = db.Count(&total)
	if offset != 0 || limit != 0 {
		db = db.Offset(offset).Limit(limit)
	}
	err = db.Find(&perms).Error
	return
}

func (r *dbRepository) GetUserPermissions(userID uint64, role Role, offset, limit int) ([]Permission, int, error) {
	var perms = make([]Permission, 0)
	var total int
	db := r.db.Model(Permission{}).
		Where("user_id = ?", userID)
	if role != 0 {
		db = db.Where("role = ?", role)
	}
	db = db.Count(&total)
	if offset != 0 || limit != 0 {
		db = db.Offset(offset).Limit(limit)
	}
	err := db.Preload("Project").Find(&perms).Error
	if err != nil {
		logger.Error(err)
		return nil, 0, errors.ProjectPermissionQueryError.Wrap(err, "cannot query project permissions for project")
	}
	return perms, total, nil
}

func (r *dbRepository) CreatePermission(projectID, userID uint64, role Role) (Permission, error) {
	perm := Permission{
		ProjectID: projectID,
		UserID:    userID,
		Role:      role,
	}
	err := r.db.Create(&perm).Error
	if err != nil {
		return perm, errors.ProjectPermissionCreatingError.Wrap(err, "cannot create project permission")
	}
	return perm, nil
}

func (r *dbRepository) UpdatePermission(projectID, userID uint64, role Role) (Permission, error) {
	var perm = Permission{
		ProjectID: projectID,
		UserID:    userID,
	}
	if r.db.Where(&perm).First(&perm).RecordNotFound() {
		return Permission{}, errors.ProjectPermissionNotFound.
			NewWithMessageF("user %d is not a member of project %d", userID, projectID)
	}
	err := r.db.Model(&perm).Update("role", role).First(&perm).Error
	if err != nil {
		return Permission{}, errors.ProjectPermissionCannotUpdate.Wrap(err, "cannot update permission")
	}
	return perm, nil
}

func (r *dbRepository) GetPermission(userID, projectID uint64) (Permission, error) {
	perm := Permission{
		UserID:    userID,
		ProjectID: projectID,
	}
	db := r.db.Where(&perm).First(&perm)
	if db.RecordNotFound() {
		return Permission{}, errors.ProjectPermissionNotFound.NewWithMessage("project permission not found")
	}
	if err := db.Error; err != nil {
		return Permission{}, errors.ProjectPermissionQueryError.NewWithMessage("error query project permission")
	}
	return perm, nil
}

func (r *dbRepository) UpdateProject(ProjectID uint64, changes map[string]interface{}) (Project, error) {
	var project Project
	project.ID = ProjectID
	err := r.db.Model(&project).Update(changes).First(&project, ProjectID).Error
	if err != nil {
		return Project{}, errors.ProjectCannotUpdate.Wrap(err, "cannot update project detail")
	}
	return project, nil
}

func (r *dbRepository) Delete(id uint64) error {
	var p Project
	p.ID = id
	err := r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Delete(&p).Error
		if err != nil {
			return err
		}
		err = tx.Where("project_id = ?", id).
			Delete(&Permission{}).Error
		if err != nil {
			return err
		}
		return nil
	})
	err = r.db.Delete(&p).Error
	if err != nil {
		return errors.ProjectCannotDelete.Wrap(err, "cannot delete project")
	}
	return nil
}

func (r *dbRepository) DeletePermission(userID, projectID uint64) error {
	var perm Permission
	perm.ProjectID = projectID
	perm.UserID = userID
	err := r.db.Model(&perm).Where(&perm).Delete(&perm).Error
	if err != nil {
		return errors.ProjectPermissionCannotDelete.NewWithMessageF("cannot delete permission for user %d, project %d", userID, projectID)
	}
	return nil
}
