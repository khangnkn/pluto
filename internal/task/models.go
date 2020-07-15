package task

import (
	"github.com/nkhang/pluto/internal/image"
	"github.com/nkhang/pluto/pkg/gorm"
	"github.com/spf13/cast"
)

type DetailStatus uint32
type Status uint32
type Role uint32

const (
	AnyRole Role = iota
	Assigner
	Labeler
	Reviewer
)

const (
	Pending DetailStatus = iota
	Draft
	Labeled
	Approved
	Rejected
)

const (
	Any Status = iota
	Labeling
	Reviewing
	Done
)

type Task struct {
	gorm.Model
	Title       string
	Description string
	ProjectID   uint64
	DatasetID   uint64
	Assigner    uint64
	Labeler     uint64
	Reviewer    uint64
	Status      Status
}

type Detail struct {
	gorm.Model
	Status  DetailStatus
	TaskID  uint64
	ImageID uint64
	Image   image.Image
}

func (d Detail) TableName() string {
	suffix := d.TaskID % 10
	return "task_detail_" + cast.ToString(suffix)
}
