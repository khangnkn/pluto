package taskapi

import (
	"github.com/nkhang/pluto/internal/dataset/datasetapi"
	"github.com/nkhang/pluto/internal/image/imageapi"
	"github.com/nkhang/pluto/internal/label/labelapi"
	"github.com/nkhang/pluto/internal/project/projectapi"
	"github.com/nkhang/pluto/internal/task"
	"github.com/nkhang/pluto/internal/workspace/workspaceapi"
)

type CreateTaskRequest struct {
	ProjectID uint64         `json:"project_id" form:"project_id"`
	DatasetID uint64         `json:"dataset_id" form:"dataset_id" binding:"required"`
	Assigner  uint64         `json:"assigner" form:"assigner" binding:"required"`
	Labeler   uint64         `json:"labeler" form:"labeler" binding:"required"`
	Quantity  int            `json:"quantity" form:"quantity" binding:"required"`
	Assignees []AssigneePair `json:"assignees" form:"assignees" binding:"required"`
}

type AssigneePair struct {
	Reviewer uint64 `json:"reviewer" form:"reviewer" binding:"required"`
	Quantity int    `json:"quantity" form:"quantity" binding:"required"`
}

type GetTaskDetailsRequest struct {
	TaskID   uint64 `form:"-"`
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
}

type UpdateTaskDetailRequest struct {
	Status uint64 `form:"status" json:"status"`
}

type TaskDetailResponse struct {
	ID     uint64                 `json:"id"`
	Status int32                  `json:"status"`
	TaskID uint64                 `json:"task_id"`
	Image  imageapi.ImageResponse `json:"image"`
}

type TaskResponse struct {
	ID          uint64 `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Assigner    uint64 `json:"assigner"`
	Labeler     uint64 `json:"labeler"`
	Reviewer    uint64 `json:"reviewer"`
	CreatedAt   uint64 `json:"created_at"`
}

type PushTaskMessage struct {
	Workspace workspaceapi.WorkspaceDetailResponse `json:"workspace"`
	Project   projectapi.ProjectResponse           `json:"project"`
	Dataset   datasetapi.DatasetResponse           `json:"dataset"`
	Tasks     []TaskResponse                       `json:"tasks"`
	Labels    []labelapi.LabelResponse             `json:"labels"`
}

func ToTaskDetailResponse(detail task.Detail) TaskDetailResponse {
	return TaskDetailResponse{
		ID:     detail.ID,
		Status: int32(detail.Status),
		TaskID: detail.TaskID,
		Image:  imageapi.ToImageResponse(detail.Image),
	}
}
