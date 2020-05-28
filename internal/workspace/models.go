package workspace

import (
	"github.com/nkhang/pluto/internal/project"
	"github.com/nkhang/pluto/pkg/gorm"
)

type Workspace struct {
	gorm.Model
	WorkspaceID uint64
	Title       string
	Description string
	Projects    []project.Project
	Perm        []Permission
}

type Permission struct {
	gorm.Model
	WorkspaceID uint64
	Workspace   Workspace `gorm:"association_save_reference:false"`
	UserID      uint64
}

func (Permission) TableName() string {
	return "workspace_permissions"
}
