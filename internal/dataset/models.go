package dataset

import (
	"github.com/nkhang/pluto/internal/project"
	"github.com/nkhang/pluto/pkg/gorm"
)

const (
	fieldProjectID = "project_id"
)

type Dataset struct {
	gorm.Model
	Title       string
	Description string
	ProjectID   uint64
	Thumbnail   string
	Project     project.Project
}
