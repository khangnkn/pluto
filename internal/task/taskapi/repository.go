package taskapi

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/nkhang/pluto/internal/workspace/workspaceapi"

	"github.com/nkhang/pluto/pkg/annotation"

	"github.com/nkhang/pluto/pkg/logger"

	"github.com/nkhang/pluto/pkg/util/clock"

	"github.com/nkhang/pluto/pkg/errors"

	"github.com/nkhang/pluto/internal/dataset/datasetapi"
	"github.com/nkhang/pluto/internal/image"
	"github.com/nkhang/pluto/internal/project/projectapi"
	"github.com/nkhang/pluto/internal/task"
	"github.com/nkhang/pluto/pkg/util/paging"
)

type Repository interface {
	GetTasks(userID uint64, request GetTasksRequest) (response GetTaskResponse, err error)
	GetTaskForProject(projectID, userID uint64, request GetTasksRequest) (response GetTaskResponse, err error)
	GetTask(taskID uint64) (TaskResponse, error)
	CreateTask(projectID, assigner uint64, request CreateTaskRequest) error
	DeleteTask(taskID uint64) error
	GetTaskDetails(taskID uint64, request GetTaskDetailsRequest) ([]TaskDetailResponse, error)
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

func (r *repository) GetTasks(userID uint64, request GetTasksRequest) (response GetTaskResponse, err error) {
	offset, limit := paging.Parse(request.Page, request.PageSize)
	var (
		tasks = make([]task.Task, 0)
		total int
	)
	switch request.Source {
	case SrcAllTasks:
		tasks, total, err = r.repository.GetTasksByUser(userID, task.AnyRole, task.Any, offset, limit)
		if err != nil {
			return GetTaskResponse{}, err
		}
	case SrcAssignerTasks:
		tasks, total, err = r.repository.GetTasksByUser(userID, task.Assigner, task.Any, offset, limit)
		if err != nil {
			return GetTaskResponse{}, err
		}
	case SrcLabelingTasks:
		tasks, total, err = r.repository.GetTasksByUser(userID, task.Labeler, task.Any, offset, limit)
		if err != nil {
			return GetTaskResponse{}, err
		}
	case SrcReviewingTasks:
		tasks, total, err = r.repository.GetTasksByUser(userID, task.Reviewer, task.Any, offset, limit)
		if err != nil {
			return GetTaskResponse{}, err
		}
	default:
		return GetTaskResponse{}, errors.TaskCannotGet.NewWithMessage("role is not supported")
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

func (r *repository) CreateTask(projectID, assigner uint64, request CreateTaskRequest) error {
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
		task, err := r.repository.CreateTask(request.Title, request.Description, assigner, pair.Labeler, pair.Reviewer, projectID, request.DatasetID, ids)
		if err != nil {
			errs = append(errs, err)
		}
		tasks = append(tasks, task)
	}
	if len(tasks) != 0 {
		err := r.annotationService.CreateTask(projectID, request.DatasetID, tasks)
		if err != nil {
			logger.Errorf("create task to annotation server error. err %v", err)
		}
	}
	if len(errs) == 0 {
		return nil
	}
	msg := fmt.Sprintf("failed to create %d tasks", len(errs))
	for i := range errs {
		logger.Infof("error creating task %v", errs[i])
	}
	return errors.TaskCannotCreate.NewWithMessage(msg)
}

func (r *repository) DeleteTask(taskID uint64) error {
	return r.repository.DeleteTask(taskID)
}

func (r *repository) GetTaskDetails(taskID uint64, request GetTaskDetailsRequest) ([]TaskDetailResponse, error) {
	details, _, err := r.repository.GetTaskDetails(taskID, request.Status, request.CurrentID, request.PageSize)
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
	err = r.repository.CheckTaskStatus(taskID, request.Status)
	if err != nil {
		logger.Error("error check update task status for task %d, detail status %d", taskID, request.Status)
	}
	if _, ok := changes["status"]; ok && request.Status == 2 {
		err := r.imgRepo.Incr(detail.ImageID)
		if err != nil {
			logger.Errorf("error increasing image status %v, id %d", err, detail.ImageID)
		}
	}
	return ToTaskDetailResponse(detail), nil
}

func (r *repository) ToTaskResponse(t task.Task) TaskResponse {
	_, total, err := r.repository.GetTaskDetails(t.ID, task.AnyStatus, 0, 0)
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
		Project: ProjectObject{
			ProjectBaseResponse: projectapi.ProjectBaseResponse{
				ID:          project.ID,
				Title:       project.Title,
				Description: project.Description,
				Thumbnail:   project.Thumbnail,
				Color:       project.Color,
			},
			ProjectManagers: project.ProjectManagers,
		},

		Workspace: WorkspaceObject{
			WorkspaceBaseResponse: workspaceapi.WorkspaceBaseResponse{
				ID:          project.Workspace.ID,
				Updated:     project.Workspace.Updated,
				Title:       project.Workspace.Title,
				Description: project.Workspace.Description,
				Color:       project.Workspace.Color,
			},
			Admin: project.Workspace.Admin,
		},
		Assigner:   t.Assigner,
		Labeler:    t.Labeler,
		Reviewer:   t.Reviewer,
		Status:     uint32(t.Status),
		ImageCount: total,
		CreatedAt:  clock.UnixMillisecondFromTime(t.CreatedAt),
	}
}

func (r *repository) GetTaskForProject(projectID, userID uint64, request GetTasksRequest) (resp GetTaskResponse, err error) {
	offset, limit := paging.Parse(request.Page, request.PageSize)
	var (
		tasks = make([]task.Task, 0)
		total int
	)
	logger.Info(request.Source)
	switch request.Source {
	case SrcAllTasks:
		tasks, total, err = r.repository.GetByProjectAndUser(projectID, userID, task.AnyRole, offset, limit)
	case SrcAssignerTasks:
		tasks, total, err = r.repository.GetByProjectAndUser(projectID, userID, task.Assigner, offset, limit)
	case SrcLabelingTasks:
		tasks, total, err = r.repository.GetByProjectAndUser(projectID, userID, task.Labeler, offset, limit)
	case SrcReviewingTasks:
		tasks, total, err = r.repository.GetByProjectAndUser(projectID, userID, task.Reviewer, offset, limit)
	default:
		return GetTaskResponse{}, errors.TaskCannotGet.NewWithMessage("role is not supported")
	}
	if err != nil {
		return
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
