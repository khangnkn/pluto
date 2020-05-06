package project

import (
	"github.com/jinzhu/gorm"

	"github.com/nkhang/pluto/pkg/errors"
)

type DiskRepository interface {
	Get(wID uint64, pID uint64) (Project, error)
	GetByWorkspaceID(wID uint64) ([]Project, error)
}

type diskRepository struct {
	db *gorm.DB
}

func NewDiskRepository(db *gorm.DB) *diskRepository {
	return &diskRepository{db: db}
}

func (r *diskRepository) Get(wID uint64, pID uint64) (Project, error) {
	var p Project
	result := r.db.Where(fieldWorkspaceID+" = ?", wID).First(&p, pID)
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
