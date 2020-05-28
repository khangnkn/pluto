package project

import (
	"github.com/nkhang/pluto/internal/label"
	"github.com/nkhang/pluto/pkg/gorm"
)

const (
	fieldWorkspaceID = "workspace_id"
)

type Role int32

const (
	Labeler Role = iota
	Reviewer
	Manager
)

type Project struct {
	gorm.Model
	Title       string
	Description string
	Labels      []label.Label
}

type Permission struct {
	gorm.Model
	ProjectID uint64
	Project   Project `gorm:"association_save_reference:false"`
	UserID    uint64
	Role      Role
}

func (Permission) TableName() string {
	return "project_permissions"
}
