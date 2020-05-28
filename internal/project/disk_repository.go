package project

import (
	"github.com/jinzhu/gorm"

	"github.com/nkhang/pluto/pkg/errors"
)

type DiskRepository interface {
	Get(pID uint64) (Project, error)
	GetByWorkspaceID(wID uint64) ([]Project, error)
	GetProjectPermission(pID uint64) ([]Permission, error)
}

type diskRepository struct {
	db *gorm.DB
}

func NewDiskRepository(db *gorm.DB) *diskRepository {
	return &diskRepository{db: db}
}

func (r *diskRepository) Get(pID uint64) (Project, error) {
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

func (r *diskRepository) GetByWorkspaceID(wID uint64) ([]Project, error) {
	var projects = make([]Project, 0)
	err := r.db.Where(fieldWorkspaceID+" = ?", wID).Find(&projects).Error
	if err != nil {
		return nil, errors.ProjectQueryError.NewWithMessage("error getting project of workspace")
	}
	return projects, nil
}

func (r *diskRepository) GetProjectPermission(pID uint64) ([]Permission, error) {
	var perms = make([]Permission, 0)
	err := r.db.Where("workspace_id = ?", pID).Find(&perms).Error
	if err != nil {
		return nil, errors.ProjectQueryError.Wrap(err, "cannot query project permissions for project")
	}
	return perms, nil
}
