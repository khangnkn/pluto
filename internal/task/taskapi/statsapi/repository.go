package statsapi

import "github.com/nkhang/pluto/internal/task"

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
	details, total, err := r.taskRepo.GetTaskDetails(taskID, task.Approved, 0, 0)
	if err != nil {
		return
	}
	resp = TaskStatsResponse{
		Processed: len(details),
		Total:     total,
	}
	return
}
