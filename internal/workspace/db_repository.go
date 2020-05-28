package workspace

import (
	"github.com/jinzhu/gorm"

	"github.com/nkhang/pluto/pkg/errors"
)

type DBRepository interface {
	Get(id uint64) (Workspace, error)
	GetByUserID(userID uint64) ([]Workspace, error)
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

func (r *dbRepository) GetByUserID(userID uint64) ([]Workspace, error) {
	var workspaces = make([]Workspace, 0)
	var perms = make([]Permission, 0)
	err := r.db.Debug().Where("user_id = ?", userID).Find(&perms).Error
	if err != nil {
		return nil, errors.WorkspaceQueryError.Wrap(err, "workspace query error")
	}
	for _, perm := range perms {
		var w Workspace
		err := r.db.Model(&perm).Association("Workspace").Find(&w).Error
		if err != nil {
			return nil, errors.WorkspaceQueryError.Wrap(err, "workspace query error")
		}
		workspaces = append(workspaces, w)
	}
	return workspaces, nil
}
