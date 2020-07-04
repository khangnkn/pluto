package taskapi

import (
	"github.com/nkhang/pluto/internal/image/imageapi"
	"github.com/nkhang/pluto/internal/task"
)

type CreateTaskRequest struct {
	DatasetID uint64 `json:"dataset_id" form:"dataset_id" binding:"required"`
	Assigner  uint64 `json:"assigner" form:"assigner" binding:"required"`
	Labeler   uint64 `json:"labeler" form:"labeler" binding:"required"`
	Reviewer  uint64 `json:"reviewer" form:"reviewer" binding:"required"`
	Quantity  int    `json:"quantity" form:"quantity" binding:"required"`
}

type GetTaskDetailsRequest struct {
	TaskID   uint64 `form:"-"`
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
}

type TaskDetailResponse struct {
	ID     uint64                 `json:"id"`
	TaskID uint64                 `json:"task_id"`
	Image  imageapi.ImageResponse `json:"image"`
}

func ToTaskDetailResponse(detail task.Detail) TaskDetailResponse {
	return TaskDetailResponse{
		ID:     detail.ID,
		TaskID: detail.TaskID,
		Image:  imageapi.ToImageResponse(detail.Image),
	}
}
