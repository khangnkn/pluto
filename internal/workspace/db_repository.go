package workspace

import (
	"github.com/jinzhu/gorm"
	"github.com/nkhang/pluto/pkg/logger"
	gormbulk "github.com/t-tiger/gorm-bulk-insert/v2"

	"github.com/nkhang/pluto/pkg/errors"
)

type DBRepository interface {
	Get(id uint64) (Workspace, error)
	GetByUserID(userID uint64, role Role, offset, limit int) ([]Workspace, int, error)
	GetPermissionByWorkspaceID(workspaceID uint64, role Role, offset, limit int) ([]Permission, int, error)
	Create(userID uint64, title, description, color string) (Workspace, error)
	DeleteWorkspace(workspaceID uint64) error
	CreatePermission(workspaceID uint64, userIDs []uint64, role Role) error
	DeletePermission(workspaceID uint64, userID uint64) error
	UpdateWorkspace(workspaceID uint64, changes map[string]interface{}) (Workspace, error)
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

func (r *dbRepository) GetByUserID(userID uint64, role Role, offset, limit int) ([]Workspace, int, error) {
	var count int
	var perms = make([]Permission, 0)
	db := r.db.Model(Permission{}).
		Where("user_id = ?", userID)
	if int32(role) != 0 {
		db = db.Where("role = ?", role)
	}
	db = db.Count(&count)
	if offset != 0 || limit != 0 {
		db = db.Offset(offset).Limit(limit)
	}
	err := db.Preload("Workspace").Find(&perms).Error
	if err != nil {
		return nil, 0, errors.WorkspaceQueryError.Wrap(err, "workspace query error")
	}
	var workspaces = make([]Workspace, len(perms))
	for i := range perms {
		workspaces[i] = perms[i].Workspace
	}
	return workspaces, count, nil
}

func (r *dbRepository) GetPermissionByWorkspaceID(workspaceID uint64, role Role, offset, limit int) ([]Permission, int, error) {
	var (
		count int
		perms = make([]Permission, 0)
	)
	db := r.db.Model(Permission{WorkspaceID: workspaceID}).
		Where("workspace_id = ?", workspaceID)
	if int32(role) != 0 {
		db = db.Where("role = ?", role)
	}
	db = db.Count(&count)
	if offset != 0 || limit != 0 {
		db = db.Offset(offset).Limit(limit)
	}
	err := db.Preload("Workspace").Find(&perms).Error
	if err != nil {
		return nil, 0, errors.WorkspaceQueryError.Wrap(err, "workspace query error")
	}
	return perms, count, nil
}

func (r *dbRepository) Create(userID uint64, title, description, color string) (Workspace, error) {
	var w = Workspace{
		Title:       title,
		Description: description,
		Color:       color,
	}
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := r.db.Save(&w).Error; err != nil {
			return errors.WorkspaceErrorCreating.Wrap(err, "cannot create workspace")
		}
		var perm = Permission{
			WorkspaceID: w.ID,
			Role:        Admin,
			UserID:      userID,
		}
		if err := r.db.Save(&perm).Error; err != nil {
			return errors.WorkspaceErrorCreating.Wrap(err, "cannot create workspace")
		}
		return nil
	})
	if err != nil {
		return Workspace{}, err
	}
	return w, nil
}

func (r *dbRepository) UpdateWorkspace(workspaceID uint64, changes map[string]interface{}) (Workspace, error) {
	var workspace Workspace
	workspace.ID = workspaceID
	db := r.db.Model(&workspace).Update(changes).First(&workspace, workspaceID)
	if db.RecordNotFound() {
		logger.Infof("is empty %v", db.RecordNotFound())
		return Workspace{}, errors.WorkspaceNotFound.NewWithMessageF("workspace %d not found", workspaceID)
	}
	if err := db.Error; err != nil {
		return Workspace{}, errors.WorkspaceCannotUpdate.Wrap(err, "cannot update workspace detail")
	}
	return workspace, nil
}

func (r *dbRepository) CreatePermission(workspaceID uint64, userIDs []uint64, role Role) error {
	var perms = make([]interface{}, len(userIDs))
	for i, userID := range userIDs {
		perms[i] = Permission{
			WorkspaceID: workspaceID,
			Role:        role,
			UserID:      userID,
		}
	}
	err := gormbulk.BulkInsert(r.db, perms, 1000)
	if err != nil {
		return errors.WorkspacePermissionErrorCreating.Wrap(err, "cannot create permissions")
	}
	return nil
}

func (r *dbRepository) DeletePermission(workspaceID uint64, userID uint64) error {
	var perm = Permission{
		WorkspaceID: workspaceID,
		UserID:      userID,
	}
	err := r.db.Delete(&perm).Error
	if err != nil {
		return errors.WorkspacePermissionDeletingError.Wrap(err, "cannot delete user from workspace")
	}
	return nil
}

func (r *dbRepository) DeleteWorkspace(workspaceID uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		db := r.db.Model(&Permission{}).Where("workspace_id = ?", workspaceID).Delete(&Permission{})
		if err := db.Error; err != nil {
			return errors.WorkspaceErrorDeleting.Wrap(err, "cannot delete workspace")
		}
		if err := r.db.Delete(&Workspace{}, workspaceID).Error; err != nil {
			return errors.WorkspaceErrorDeleting.Wrap(err, "cannot delete workspace")
		}
		return nil
	})
}
