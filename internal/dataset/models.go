package dataset

import (
	"github.com/nkhang/pluto/pkg/gorm"
)

const (
	fieldProjectID = "project_id"
)

type Dataset struct {
	gorm.Model
	Title       string
	Description string
	Thumbnail   string
	ProjectID   uint64
}
