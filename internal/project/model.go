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
	Any Role = iota
	Manager
	Member
	Admin
)

var defaultImage = "http://annotation.ml:9000/plutos3/placeholder.png"

type Project struct {
	gorm.Model
	WorkspaceID uint64
	Title       string
	Description string
	Thumbnail   string
	Color       string
	Dir         string
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
