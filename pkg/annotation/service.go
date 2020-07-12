package annotation

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/spf13/viper"

	"github.com/nkhang/pluto/internal/dataset"
	"github.com/nkhang/pluto/internal/label"
	"github.com/nkhang/pluto/internal/project"
	"github.com/nkhang/pluto/internal/task"
	"github.com/nkhang/pluto/internal/workspace"
	"github.com/nkhang/pluto/pkg/errors"
	"github.com/nkhang/pluto/pkg/logger"
	"github.com/nkhang/pluto/pkg/util/clock"
)

type Service interface {
	CreateTask(projectID, datasetID uint64, tasks []task.Task) error
}

type service struct {
	workspaceRepo      workspace.Repository
	projectRepo        project.Repository
	datasetRepo        dataset.Repository
	labelRepo          label.Repository
	client             http.Client
	annotationBasePath string
}

func NewService(workspaceRepo workspace.Repository,
	projectRepo project.Repository,
	datasetRepo dataset.Repository,
	labelRepo label.Repository) *service {
	client := http.Client{}
	annotationBase := viper.GetString("annotation.baseurl")
	return &service{
		workspaceRepo:      workspaceRepo,
		projectRepo:        projectRepo,
		datasetRepo:        datasetRepo,
		labelRepo:          labelRepo,
		client:             client,
		annotationBasePath: annotationBase,
	}
}

func (s *service) CreateTask(projectID, datasetID uint64, tasks []task.Task) error {
	p, err := s.projectRepo.Get(projectID)
	if err != nil {
		return err
	}
	message, err := NewBuilder(
		s.workspaceRepo,
		s.projectRepo,
		s.datasetRepo,
		s.labelRepo).
		WithWorkspace(p.WorkspaceID).
		WithProject(projectID).
		WithDataset(datasetID).
		WithTasks(tasks).
		WithLabels(projectID).
		Build()
	if err != nil {
		return err
	}
	return s.push(message)
}

func (s *service) push(message PushTaskMessage) error {
	path := s.annotationBasePath + "/task"
	b, err := json.Marshal(&message)
	if err != nil {
		return err
	}

	resp, err := s.client.Post(path, "application/json", bytes.NewReader(b))
	if err != nil {
		return err
	}
	bb, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		logger.Info(string(bb))
	}
	return nil
}

type builder struct {
	errs          []error
	workspaceRepo workspace.Repository
	projectRepo   project.Repository
	datasetRepo   dataset.Repository
	labelRepo     label.Repository

	workspace WorkspaceObject
	project   ProjectObject
	dataset   DatasetObject
	tasks     []TaskObject
	labels    []LabelObject
}

func NewBuilder(workspaceRepo workspace.Repository,
	projectRepo project.Repository,
	datasetRepo dataset.Repository,
	labelRepo label.Repository) *builder {
	return &builder{
		errs:          make([]error, 0),
		workspaceRepo: workspaceRepo,
		projectRepo:   projectRepo,
		datasetRepo:   datasetRepo,
		labelRepo:     labelRepo,
	}
}

func (b *builder) WithWorkspace(id uint64) *builder {
	w, err := b.workspaceRepo.Get(id)
	if err != nil {
		b.errs = append(b.errs, err)
		return b
	}
	var admin uint64
	perms, _, err := b.workspaceRepo.GetPermission(w.ID, workspace.Admin, 0, 1)
	if err != nil || len(perms) == 0 {
		b.errs = append(b.errs, errors.WorkspacePermissionErrorCreating.NewWithMessage(""))
	} else {
		admin = perms[0].UserID
	}
	b.workspace = WorkspaceObject{
		ID:    w.ID,
		Title: w.Title,
		Admin: admin,
	}
	return b
}

func (b *builder) WithProject(id uint64) *builder {
	p, err := b.projectRepo.Get(id)
	if err != nil {
		b.errs = append(b.errs, err)
		return b
	}
	var manager uint64
	perms, _, err := b.projectRepo.GetProjectPermissions(p.ID, project.Manager, 0, 1)
	if err != nil || len(perms) == 0 {
		b.errs = append(b.errs, errors.WorkspacePermissionErrorCreating.NewWithMessage(""))
		return b
	} else {
		manager = perms[0].UserID
	}
	b.project = ProjectObject{
		ID:             p.ID,
		Title:          p.Title,
		ProjectManager: manager,
	}
	return b
}

func (b *builder) WithDataset(id uint64) *builder {
	d, err := b.datasetRepo.Get(id)
	if err != nil {
		b.errs = append(b.errs, err)
		return b
	}
	b.dataset = DatasetObject{
		ID:        d.ID,
		Title:     d.Title,
		ProjectID: d.ProjectID,
	}
	return b
}

func (b *builder) WithTasks(tasks []task.Task) *builder {
	var t = make([]TaskObject, len(tasks))
	for i, task := range tasks {
		t[i] = TaskObject{
			ID:        task.ID,
			Labeler:   task.Labeler,
			Reviewer:  task.Reviewer,
			CreatedAt: clock.UnixMillisecondFromTime(task.CreatedAt),
		}
	}
	b.tasks = t
	return b
}

func (b *builder) WithLabels(projectID uint64) *builder {
	labels, err := b.labelRepo.GetByProjectId(projectID)
	if err != nil {
		b.errs = append(b.errs, err)
		return b
	}
	var responses = make([]LabelObject, len(labels))
	for i, label := range labels {
		responses[i] = LabelObject{
			ID:    label.ID,
			Name:  label.Name,
			Color: label.Color,
			Tool: ToolObject{
				ID:   label.Tool.ID,
				Name: label.Tool.Name,
			},
		}
	}
	b.labels = responses
	return b
}
func (b *builder) Build() (PushTaskMessage, error) {
	if len(b.errs) != 0 {
		logger.Errorf("error creating task %v", b.errs)
		return PushTaskMessage{}, errors.TaskCannotCreate.NewWithMessage("cannot build push task message")
	}
	return PushTaskMessage{
		Workspace: b.workspace,
		Project:   b.project,
		Dataset:   b.dataset,
		Tasks:     b.tasks,
		Labels:    b.labels,
	}, nil
}
