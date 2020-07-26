package statsapi

import (
	"github.com/nkhang/pluto/internal/task"
)

type Repository interface {
	Stats(taskID uint64) (resp TaskStatsResponse, err error)
}

type repository struct {
	taskRepo task.Repository
}

func NewRepository(r task.Repository) *repository {
	return &repository{taskRepo: r}
}

func (r *repository) Stats(taskID uint64) (resp TaskStatsResponse, err error) {
	t, err := r.taskRepo.GetTask(taskID)
	if err != nil {
		return
	}
	switch t.Status {
	case task.Labeling:
		return r.buildForLabeling(t)
	default:
		return r.buildForReviewing(t)
	}
}

func (r *repository) buildForLabeling(d task.Task) (resp TaskStatsResponse, err error) {
	details, total, err := r.taskRepo.GetTaskDetails(d.ID, task.Labeled, 0, 0)
	if err != nil {
		return
	}
	resp = TaskStatsResponse{
		Processed: len(details),
		Total:     total,
	}
	return
}

func (r *repository) buildForReviewing(d task.Task) (resp TaskStatsResponse, err error) {
	d1, total, err := r.taskRepo.GetTaskDetails(d.ID, task.Approved, 0, 0)
	if err != nil {
		return
	}
	d2, _, err := r.taskRepo.GetTaskDetails(d.ID, task.Rejected, 0, 0)
	if err != nil {
		return
	}
	resp = TaskStatsResponse{
		Processed: len(append(d1, d2...)),
		Total:     total,
	}
	return
}
