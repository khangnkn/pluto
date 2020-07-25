package dataset

import (
	"github.com/nkhang/pluto/pkg/gorm"
)

const (
	fieldProjectID = "project_id"
	defaultImage   = "http://annotation.ml:9000/plutos3/placeholder.png"
)

type Dataset struct {
	gorm.Model
	Title       string
	Description string
	Thumbnail   string
	ProjectID   uint64
}
