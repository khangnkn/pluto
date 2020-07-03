package task

import (
	"github.com/nkhang/pluto/pkg/gorm"
	"github.com/spf13/cast"
)

type DetailStatus uint32

const (
	Unassigned DetailStatus = iota
	Assigned
	Labeled
	Reviewed
)

var DetailStatusMap = map[DetailStatus]string{
	Unassigned: "Unassigned",
	Assigned:   "Assigned",
	Labeled:    "Labeled",
	Reviewed:   "Reviewed",
}

type Task struct {
	gorm.Model
	DatasetID uint64
	Assigner  uint64
	Labeler   uint64
	Reviewer  uint64
}

type Detail struct {
	gorm.Model
	Status  DetailStatus
	TaskID  uint64
	ImageID uint64
}

func (d Detail) TableName() string {
	suffix := d.TaskID % 10
	return "task_detail_" + cast.ToString(suffix)
}
