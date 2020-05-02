package workspace

import (
	"github.com/jinzhu/gorm"

	"github.com/nkhang/pluto/pkg/errors"
)

type DBRepository interface {
	Get(id uint64) (Workspace, error)
}

type dbRepository struct {
	db *gorm.DB
}

func NewDDBRepository(db *gorm.DB) *dbRepository {
	return &dbRepository{db: db}
}

func (r *dbRepository) Get(id uint64) (Workspace, error) {
	var w Workspace
	result := r.db.First(&w, id)
	if result.RecordNotFound() {
		return Workspace{}, errors.WorkspaceNotFound.NewWithMessage("workspace not found")

	}
	if err := result.Error; err != nil {
		return Workspace{}, errors.WorkspaceQueryError.Wrap(err, "workspace query error")
	}
	return w, nil
}
