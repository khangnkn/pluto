package taskapi

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/nkhang/pluto/pkg/annotation"

	"github.com/nkhang/pluto/pkg/logger"

	"github.com/nkhang/pluto/pkg/util/clock"

	"github.com/nkhang/pluto/pkg/errors"

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
	GetTasks(request GetTasksRequest) (response GetTaskResponse, err error)
	GetTask(taskID uint64) (TaskResponse, error)
	CreateTask(request CreateTaskRequest) error
	DeleteTask(taskID uint64) error
	GetTaskDetails(request GetTaskDetailsRequest) ([]TaskDetailResponse, error)
	UpdateTaskDetail(taskID, detailID uint64, request UpdateTaskDetailRequest) (TaskDetailResponse, error)
}

type repository struct {
	repository        task.Repository
	imgRepo           image.Repository
	datasetRepo       datasetapi.Repository
	projectRepo       projectapi.Repository
	annotationService annotation.Service
}

func NewRepository(r task.Repository,
	ir image.Repository,
	datasetRepo datasetapi.Repository,
	projectRepo projectapi.Repository,
	annotationService annotation.Service) *repository {
	return &repository{
		repository:        r,
		imgRepo:           ir,
		datasetRepo:       datasetRepo,
		projectRepo:       projectRepo,
		annotationService: annotationService,
	}
}

func (r *repository) GetTask(taskID uint64) (TaskResponse, error) {
	task, err := r.repository.GetTask(taskID)
	if err != nil {
		return TaskResponse{}, err
	}
	return r.ToTaskResponse(task), nil
}

func (r *repository) GetTasks(request GetTasksRequest) (response GetTaskResponse, err error) {
	offset, limit := paging.Parse(request.Page, request.PageSize)
	var (
		tasks = make([]task.Task, 0)
		total int
	)
	switch request.Source {
	case SrcAllTasks:
		tasks, total, err = r.repository.GetTasksByUser(request.UserID, task.AnyRole, task.Any, offset, limit)
		if err != nil {
			return GetTaskResponse{}, err
		}
	case SrcLabelingTasks:
		tasks, total, err = r.repository.GetTasksByUser(request.UserID, task.Labeler, task.Any, offset, limit)
		if err != nil {
			return GetTaskResponse{}, err
		}
	case SrcReviewingTasks:
		tasks, total, err = r.repository.GetTasksByUser(request.UserID, task.Reviewer, task.Any, offset, limit)
		if err != nil {
			return GetTaskResponse{}, err
		}
	case SrcProjectTasks:
		if request.ProjectID == 0 {
			return GetTaskResponse{}, errors.BadRequest.NewWithMessage("project_id is required")
		}
		tasks, total, err = r.repository.GetTasksByProject(request.ProjectID, task.Any, offset, limit)
		if err != nil {
			return GetTaskResponse{}, err
		}
	}
	responses := make([]TaskResponse, len(tasks))
	for i := range tasks {
		responses[i] = r.ToTaskResponse(tasks[i])
	}
	return GetTaskResponse{
		Total: total,
		Tasks: responses,
	}, nil
}

func (r *repository) CreateTask(request CreateTaskRequest) error {
	imgs, err := r.imgRepo.GetAllImageByDataset(request.DatasetID)
	if err != nil {
		return err
	}
	if len(imgs) == 0 {
		return errors.TaskCannotCreate.NewWithMessageF("dataset %d has no images, abort", request.DatasetID)
	}
	var errs = make([]error, 0)
	var cursor = 0
	var tasks = make([]task.Task, 0)
	for _, pair := range request.Assignees {
		truncated := truncate(imgs, &cursor, request.Quantity)
		ids := make([]uint64, len(truncated))
		for j := range truncated {
			ids[j] = truncated[j].ID
		}
		task, err := r.repository.CreateTask(request.Title, request.Description, request.Assigner, pair.Labeler, pair.Reviewer, request.ProjectID, request.DatasetID, ids)
		if err != nil {
			errs = append(errs, err)
		}
		tasks = append(tasks, task)
	}
	if len(tasks) != 0 {
		err := r.annotationService.CreateTask(request.ProjectID, request.DatasetID, tasks)
		if err != nil {
			logger.Error(err)
		}
	}
	if len(errs) == 0 {
		return nil
	}
	msg := fmt.Sprintf("failed to create %d tasks", len(errs))
	return errors.TaskCannotCreate.NewWithMessage(msg)
}

func (r *repository) DeleteTask(taskID uint64) error {
	return r.repository.DeleteTask(taskID)
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

func (r *repository) ToTaskResponse(t task.Task) TaskResponse {
	var imageCount int
	details, err := r.repository.GetTaskDetails(t.ID, 0, 0)
	if err == nil {
		imageCount = len(details)
	}
	dataset, err := r.datasetRepo.GetByID(t.DatasetID)
	if err != nil {
		logger.Errorf("cannot get dataset response. error %v", err)
	}
	project, err := r.projectRepo.GetByID(t.ProjectID)
	if err != nil {
		logger.Errorf("cannot get project response. error %v", err)
	}
	return TaskResponse{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		Dataset:     dataset,
		Project:     project.Title,
		Workspace:   project.Workspace.Title,
		Assigner:    t.Assigner,
		Labeler:     t.Labeler,
		Reviewer:    t.Reviewer,
		Status:      uint32(t.Status),
		ImageCount:  imageCount,
		CreatedAt:   clock.UnixMillisecondFromTime(t.CreatedAt),
	}
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

func truncate(imgs []image.Image, cursor *int, s int) (res []image.Image) {
	l := len(imgs)
	if s >= l {
		*cursor = 0
		return imgs
	}
	if position := *cursor + s; position <= l {
		log.Print(*cursor, position)
		res = imgs[*cursor:position]
		*cursor += s
	} else {
		left := s - (l - *cursor)
		log.Print(left)
		res = append(imgs[*cursor:l], imgs[:left]...)
		*cursor = left
	}
	return
}
