package projectapi

import (
	"encoding/json"

	"github.com/nkhang/pluto/pkg/util/clock"

	"github.com/nkhang/pluto/internal/workspace/workspaceapi"

	"github.com/nkhang/pluto/internal/dataset"
	"github.com/nkhang/pluto/internal/project"
	"github.com/nkhang/pluto/pkg/errors"
	"github.com/nkhang/pluto/pkg/logger"
	"github.com/nkhang/pluto/pkg/util/paging"
)

type Repository interface {
	GetByID(pID uint64) (ProjectResponse, error)
	GetList(p GetProjectRequest) ([]ProjectResponse, int, error)
	GetPermissions(projectID uint64) (PermissionResponse, error)
	Create(p CreateProjectRequest) (ProjectResponse, error)
	CreatePerm(p CreatePermParams) error
	UpdateProject(id uint64, request UpdateProjectRequest) (ProjectResponse, error)
	UpdatePerm(projectID uint64, req UpdatePermissionRequest) (PermissionObject, error)
	DeleteProject(id uint64) error
}

type repository struct {
	repository    project.Repository
	datasetRepo   dataset.Repository
	workspaceRepo workspaceapi.Repository
}

func NewRepository(r project.Repository, dr dataset.Repository, wr workspaceapi.Repository) *repository {
	return &repository{
		repository:    r,
		datasetRepo:   dr,
		workspaceRepo: wr,
	}
}

func (r *repository) GetByID(pID uint64) (ProjectResponse, error) {
	p, err := r.repository.Get(pID)
	if err != nil {
		return ProjectResponse{}, err
	}
	return r.convertResponse(p), nil
}

func (r *repository) GetList(p GetProjectRequest) (responses []ProjectResponse, total int, err error) {
	offset, limit := paging.Parse(p.Page, p.PageSize)
	var projects []project.Project
	switch p.Source {
	case SrcAllProjectInWorkspace:
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
	case SrcAllProject:
		var perms []project.Permission
		perms, total, err = r.repository.GetUserPermissions(p.UserID, project.Any, offset, limit)
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

func (r *repository) GetPermissions(projectID uint64) (PermissionResponse, error) {
	perms, total, err := r.repository.GetProjectPermissions(projectID, project.Any, 0, 0)
	if err != nil {
		return PermissionResponse{}, err
	}
	var responses = make([]PermissionObject, len(perms))
	for i := range perms {
		responses[i] = convertPermissionObject(perms[i])
	}
	return PermissionResponse{
		Total:   total,
		Members: responses,
	}, nil
}

func (r *repository) Create(p CreateProjectRequest) (ProjectResponse, error) {
	project, err := r.repository.CreateProject(p.WorkspaceID, p.Title, p.Description, p.Color)
	if err != nil {
		return ProjectResponse{}, err
	}
	return r.convertResponse(project), nil
}

func (r *repository) CreatePerm(req CreatePermParams) error {
	var errs = make([]error, 0)
	for _, p := range req.Members {
		_, err := r.repository.GetPermission(p.UserID, req.ProjectID)
		if err == nil {
			continue
		}
		_, err = r.repository.CreatePermission(req.ProjectID, p.UserID, project.Role(p.Role))
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) != 0 {
		return errors.ProjectPermissionCreatingError.NewWithMessageF("error creating %d permissions", len(errs))
	}
	return nil
}

func (r *repository) UpdatePerm(projectID uint64, req UpdatePermissionRequest) (PermissionObject, error) {
	var role project.Role
	switch req.Role {
	case 1:
		role = project.Manager
	case 2:
		role = project.Member
	default:
		return PermissionObject{}, errors.ProjectRoleInvalid.NewWithMessageF("role %d is not supported", req.Role)

	}
	perm, err := r.repository.UpdatePermission(projectID, req.UserID, role)
	if err != nil {
		return PermissionObject{}, err
	}
	return convertPermissionObject(perm), nil
}

func (r *repository) UpdateProject(id uint64, request UpdateProjectRequest) (ProjectResponse, error) {
	var changes = make(map[string]interface{})
	b, _ := json.Marshal(&request)
	_ = json.Unmarshal(b, &changes)
	project, err := r.repository.UpdateProject(id, changes)
	if err != nil {
		return ProjectResponse{}, nil
	}
	return r.convertResponse(project), nil
}

func (r *repository) DeleteProject(id uint64) error {
	return r.repository.Delete(id)
}

func (r *repository) convertResponse(p project.Project) ProjectResponse {
	var datasetCount int
	d, err := r.datasetRepo.GetByProject(p.ID)
	if err != nil {
		logger.Error("error getting dataset by project id")
	} else {
		datasetCount = len(d)
	}
	var pm uint64
	perms, totalPerms, err := r.repository.GetProjectPermissions(p.ID, project.Any, 0, 0)
	if err != nil {
		logger.Error("error getting project perm")
	}
	for i := range perms {
		if perms[i].Role == project.Manager {
			pm = perms[i].UserID
			break
		}
	}
	w, _ := r.workspaceRepo.GetByID(p.WorkspaceID)
	return ProjectResponse{
		ID:             p.ID,
		Title:          p.Title,
		Description:    p.Description,
		Thumbnail:      p.Thumbnail,
		Color:          p.Color,
		DatasetCount:   datasetCount,
		MemberCount:    totalPerms,
		Workspace:      w,
		ProjectManager: pm,
	}
}

func convertPermissionObject(perm project.Permission) PermissionObject {
	return PermissionObject{
		CreatedAt: clock.UnixMillisecondFromTime(perm.CreatedAt),
		UserID:    perm.UserID,
		Role:      perm.Role,
	}
}
