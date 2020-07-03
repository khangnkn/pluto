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
	Manager Role = iota + 1
	Member
)

var Color = []string{
	"#4773AA",
	"#E4888B",
	"#7CB287",
	"#FBDB88",
	"#5FC7E3",
}

type Project struct {
	gorm.Model
	WorkspaceID uint64
	Title       string
	Description string
	Color       string
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
