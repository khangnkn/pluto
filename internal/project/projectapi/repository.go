package projectapi

import (
	"encoding/json"

	"github.com/nkhang/pluto/internal/dataset"
	"github.com/nkhang/pluto/internal/project"
	"github.com/nkhang/pluto/pkg/errors"
	"github.com/nkhang/pluto/pkg/logger"
	"github.com/nkhang/pluto/pkg/util/paging"
)

type Repository interface {
	GetByID(pID uint64) (ProjectResponse, error)
	GetList(p GetProjectParam) ([]ProjectResponse, int, error)
	Create(p CreateProjectParams) error
	CreatePerm(p CreatePermParams) error
	UpdateProject(id uint64, request UpdateProjectRequest) (ProjectResponse, error)
	DeleteProject(id uint64) error
}

type repository struct {
	repository  project.Repository
	datasetRepo dataset.Repository
}

func NewRepository(r project.Repository, dr dataset.Repository) *repository {
	return &repository{
		repository:  r,
		datasetRepo: dr,
	}
}

func (r *repository) GetByID(pID uint64) (ProjectResponse, error) {
	p, err := r.repository.Get(pID)
	if err != nil {
		return ProjectResponse{}, err
	}
	return r.convertResponse(p), nil
}

func (r *repository) GetList(p GetProjectParam) (responses []ProjectResponse, total int, err error) {
	offset, limit := paging.Parse(p.Page, p.PageSize)
	var projects []project.Project
	switch p.Source {
	case SrcAllProject:
		projects, total, err = r.repository.GetByWorkspaceID(p.WorkspaceID, offset, limit)
	case SrcMyProject:
		var perms []project.Permission
		perms, total, err = r.repository.GetUserPermissions(p.UserID, project.Manager, offset, limit)
		for i := range perms {
			projects = append(projects, perms[i].Project)
		}
	case SrcOtherProject:
		var perms []project.Permission
		perms, total, err = r.repository.GetUserPermissions(p.UserID, project.Member, offset, limit)
		for i := range perms {
			projects = append(projects, perms[i].Project)
		}
	default:
		return nil, 0, errors.BadRequest.NewWithMessage("invalid src params")
	}
	if err != nil {
		return nil, 0, err
	}
	responses = make([]ProjectResponse, len(projects))
	for i := range projects {
		responses[i] = r.convertResponse(projects[i])
	}
	return responses, total, nil
}

func (r *repository) Create(p CreateProjectParams) error {
	_, err := r.repository.CreateProject(p.WorkspaceID, p.Title, p.Desc)
	if err != nil {
		return err
	}
	go func() {
		err := r.repository.InvalidateProjectsByWorkspaceID(p.WorkspaceID)
		if err != nil {
			logger.Errorf("error invalidate cache for projects by workspace %d, error %v", p.WorkspaceID, err)
			return
		}
		logger.Infof("invalidate cache for projects by workspace %d successfully", p.WorkspaceID)
	}()
	return nil
}

func (r *repository) convertResponse(p project.Project) ProjectResponse {
	var datasetCount int
	var memberCount int
	d, err := r.datasetRepo.GetByProject(p.ID)
	if err != nil {
		logger.Error("error getting project by project id")
	} else {
		datasetCount = len(d)
	}
	m, err := r.repository.GetProjectPermissions(p.ID)
	if err != nil {
		logger.Error("error getting project perm")
	} else {
		memberCount = len(m)
	}
	return ProjectResponse{
		ID:           p.ID,
		Title:        p.Title,
		Description:  p.Description,
		Thumbnail:    p.Thumbnail,
		Color:        p.Color,
		DatasetCount: datasetCount,
		MemberCount:  memberCount,
	}
}

func (r *repository) CreatePerm(p CreatePermParams) error {
	_, err := r.repository.GetPermission(p.UserID, p.ProjectID)
	if err == nil {
		return errors.ProjectPermissionExisted.NewWithMessage("user existed in dataset")
	}
	_, err = r.repository.CreatePermission(p.ProjectID, p.UserID, project.Role(p.Role))
	if err != nil {
		return err
	}
	go func() {
		err := r.repository.InvalidatePermissionForUser(p.UserID)
		if err != nil {
			logger.Errorf("error invalidate permission for user %d", p.UserID)
		} else {
			logger.Infof("invalidate project permission for user %d successfully", p.UserID)
		}
	}()
	return nil
}

func (r *repository) UpdateProject(id uint64, request UpdateProjectRequest) (ProjectResponse, error) {
	var changes = make(map[string]interface{})
	b, _ := json.Marshal(&request)
	_ = json.Unmarshal(b, &changes)
	w, err := r.repository.UpdateProject(id, changes)
	if err != nil {
		return ProjectResponse{}, nil
	}
	return r.convertResponse(w), nil
}

func (r *repository) DeleteProject(id uint64) error {
	return r.repository.Delete(id)
}
