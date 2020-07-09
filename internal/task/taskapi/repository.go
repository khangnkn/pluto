package taskapi

import (
	"encoding/json"

	"github.com/nkhang/pluto/internal/workspace/workspaceapi"

	"github.com/nkhang/pluto/internal/dataset/datasetapi"
	"github.com/nkhang/pluto/internal/image"
	"github.com/nkhang/pluto/internal/label/labelapi"
	"github.com/nkhang/pluto/internal/project/projectapi"
	"github.com/nkhang/pluto/internal/task"
	"github.com/nkhang/pluto/internal/tool/toolapi"
	"github.com/nkhang/pluto/pkg/util/paging"
)

type Repository interface {
	CreateTask(request CreateTaskRequest) error
	GetTaskDetails(request GetTaskDetailsRequest) ([]TaskDetailResponse, error)
	UpdateTaskDetail(taskID, detailID uint64, request UpdateTaskDetailRequest) (TaskDetailResponse, error)
}

type repository struct {
	repository task.Repository
	imgRepo    image.Repository
}

func NewRepository(r task.Repository, ir image.Repository) *repository {
	return &repository{
		repository: r,
		imgRepo:    ir,
	}
}

func (r *repository) CreateTask(request CreateTaskRequest) error {
	imgs, err := r.imgRepo.GetAllImageByDataset(request.DatasetID)
	var cursor = 0
	if err != nil {
		return err
	}
	truncated := truncate(imgs, &cursor, request.Quantity)
	ids := make([]uint64, len(truncated))
	for i := range truncated {
		ids[i] = truncated[i].ID
	}
	return r.repository.CreateTask(request.Assigner, request.Labeler, 2342, request.DatasetID, ids)
}

func (r *repository) GetTaskDetails(request GetTaskDetailsRequest) ([]TaskDetailResponse, error) {
	offset, limit := paging.Parse(request.Page, request.PageSize)
	details, err := r.repository.GetTaskDetails(request.TaskID, offset, limit)
	if err != nil {
		return nil, err
	}
	var responses = make([]TaskDetailResponse, len(details))
	for i := range details {
		responses[i] = ToTaskDetailResponse(details[i])
	}
	return responses, nil
}

func (r *repository) UpdateTaskDetail(taskID, detailID uint64, request UpdateTaskDetailRequest) (TaskDetailResponse, error) {
	var changes = make(map[string]interface{})
	b, _ := json.Marshal(&request)
	_ = json.Unmarshal(b, &changes)
	detail, err := r.repository.UpdateTaskDetail(taskID, detailID, changes)
	if err != nil {
		return TaskDetailResponse{}, err
	}
	return ToTaskDetailResponse(detail), nil

}

func pushTaskMessage() PushTaskMessage {
	msg := PushTaskMessage{
		Workspace: workspaceapi.WorkspaceDetailResponse{},
		Project: projectapi.ProjectResponse{
			ID:           434,
			Title:        "fdsfsdf",
			Description:  "dfsf",
			Thumbnail:    "fsdf",
			Color:        "fsdfdfsd",
			DatasetCount: 6,
			MemberCount:  4,
		},
		Dataset: datasetapi.DatasetResponse{
			ID:          2342,
			Title:       "dsfs",
			Description: "fsf",
			ProjectID:   5,
		},
		Tasks: []TaskResponse{
			TaskResponse{
				Title:       "sdfasdf",
				Description: "sdfasdf",
				Assigner:    2343244,
				Labeler:     024,
				Reviewer:    04234,
				CreatedAt:   064,
			},
		},
		Labels: []labelapi.LabelResponse{
			labelapi.LabelResponse{
				ID:    45634,
				Name:  "ff",
				Color: "dgsdfg",
				Tool: toolapi.ToolResponse{
					ID:   1,
					Name: "rect",
				},
			},
		},
	}
	return msg
}

func truncate(imgs []image.Image, cursor *int, s int) []image.Image {
	l := len(imgs)
	if position := *cursor + s; position <= l {
		*cursor += s
		return imgs[*cursor:position]
	} else {
		left := s - (l - *cursor)
		*cursor = left
		return imgs[:left]
	}
}
