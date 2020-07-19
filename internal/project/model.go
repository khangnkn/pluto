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
	Admin
	Manager
	Member
)

var defaultImage = "https://media3.s-nbcnews.com/j/newscms/2019_33/2203981/171026-better-coffee-boost-se-329p_67dfb6820f7d3898b5486975903c2e51.fit-760w.jpg"

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
