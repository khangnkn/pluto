package image

import "github.com/nkhang/pluto/pkg/gorm"

type Status uint32

const (
	Unassigned Status = iota
	Assigned
	Labeled
	Reviewed
)

type Image struct {
	gorm.Model
	URL       string
	Status    Status
	Title     string
	Width     int
	Height    int
	Size      int64
	DatasetID uint64
}
