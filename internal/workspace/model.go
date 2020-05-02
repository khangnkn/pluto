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
}
