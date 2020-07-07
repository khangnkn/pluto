package workspace

import (
	"github.com/nkhang/pluto/pkg/gorm"
)

type Role int32

const (
	Any Role = iota
	Admin
	Member
)

type Workspace struct {
	gorm.Model
	Title       string
	Description string
}

type Permission struct {
	gorm.Model
	WorkspaceID uint64
	Workspace   Workspace `gorm:"association_save_reference:false"`
	Role        Role
	UserID      uint64
}

func (Permission) TableName() string {
	return "workspace_permissions"
}
