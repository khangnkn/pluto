package project

import (
	"github.com/nkhang/pluto/internal/label"
	"github.com/nkhang/pluto/pkg/gorm"
)

const (
	fieldWorkspaceID = "workspace_id"
)

type Project struct {
	gorm.Model
	Title       string
	Description string
	Labels      []label.Label
}
