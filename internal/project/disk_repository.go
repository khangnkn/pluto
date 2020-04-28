package project

import (
	"github.com/jinzhu/gorm"

	"github.com/nkhang/pluto/pkg/errors"
)

type DiskRepository interface {
	Get(id uint64) (Project, error)
}

type diskRepository struct {
	db *gorm.DB
}

func NewDiskRepository(db *gorm.DB) *diskRepository {
	return &diskRepository{db: db}
}

func (r *diskRepository) Get(id uint64) (Project, error) {
	var p Project
	result := r.db.First(&p, id)
	if result.RecordNotFound() {
		return Project{}, errors.ProjectNotFound.NewWithMessage("project not found")
	}
	if err := result.Error; err != nil {
		return Project{}, errors.ProjectQueryError.Wrap(err, "query project error")
	}
	return p, nil
}
