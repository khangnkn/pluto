package projectapi

import (
	"encoding/json"

	"github.com/nkhang/pluto/internal/workspace/workspaceapi"

	"github.com/nkhang/pluto/internal/dataset"
	"github.com/nkhang/pluto/internal/project"
	"github.com/nkhang/pluto/pkg/errors"
	"github.com/nkhang/pluto/pkg/logger"
	"github.com/nkhang/pluto/pkg/util/paging"
)

type Repository interface {
	GetByID(pID uint64) (ProjectResponse, error)
	GetList(userID uint64, p GetProjectRequest) ([]ProjectResponse, int, error)
	GetForWorkspace(workspaceID uint64, paging paging.Paging) (GetProjectResponse, error)
	Create(workspaceID, creator uint64, p CreateProjectRequest) (ProjectResponse, error)
	UpdateProject(id uint64, request UpdateProjectRequest) (ProjectResponse, error)
	DeleteProject(id uint64) error
	ConvertResponse(p project.Project) ProjectResponse
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
	return r.ConvertResponse(p), nil
}

func (r *repository) GetList(userID uint64, p GetProjectRequest) (responses []ProjectResponse, total int, err error) {
	offset, limit := paging.Parse(p.Page, p.PageSize)
	var projects []project.Project
	switch p.Source {
	case SrcAllProject:
		var perms []project.Permission
		perms, total, err = r.repository.GetUserPermissions(userID, project.Any, offset, limit)
		for i := range perms {
			projects = append(projects, perms[i].Project)
		}
	case SrcMyProject:
		var perms []project.Permission
		perms, total, err = r.repository.GetUserPermissions(userID, project.Manager, offset, limit)
		for i := range perms {
			projects = append(projects, perms[i].Project)
		}
	case SrcOtherProject:
		var perms []project.Permission
		perms, total, err = r.repository.GetUserPermissions(userID, project.Member, offset, limit)
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
		responses[i] = r.ConvertResponse(projects[i])
	}
	return responses, total, nil
}

func (r *repository) GetForWorkspace(workspaceID uint64, userID uint64, paging paging.Paging) (resp GetProjectResponse, err error) {
	offset, limit := paging.Parse()
	prj := make([]project.Project, 0)
	perms, _, err := r.repository.GetUserPermissions(userID, project.Any, 0, 0)
	if err != nil {
		return
	}
	for _, p := range perms {
		if p.Project.WorkspaceID == workspaceID {
			prj = append(prj, p.Project)
		}
	}
	projects := slice(prj, offset, limit)
	responses := make([]ProjectResponse, len(projects))
	for i := range projects {
		responses[i] = r.ConvertResponse(projects[i])
	}
	return GetProjectResponse{
		Total:    len(perms),
		Projects: responses,
	}, nil
}

func (r *repository) Create(workspaceID, creator uint64, p CreateProjectRequest) (ProjectResponse, error) {
	prj, err := r.repository.CreateProject(workspaceID, p.Title, p.Description, p.Color)
	if err != nil {
		return ProjectResponse{}, err
	}
	_, err = r.repository.CreatePermission(prj.ID, creator, project.Admin)
	if err != nil {
		logger.Errorf("error create admin permission for user %d to project %d workspace %d", creator, prj.ID, workspaceID)
	}
	return r.ConvertResponse(prj), nil
}

func (r *repository) UpdateProject(id uint64, request UpdateProjectRequest) (ProjectResponse, error) {
	var changes = make(map[string]interface{})
	b, _ := json.Marshal(&request)
	_ = json.Unmarshal(b, &changes)
	project, err := r.repository.UpdateProject(id, changes)
	if err != nil {
		return ProjectResponse{}, nil
	}
	return r.ConvertResponse(project), nil
}

func (r *repository) DeleteProject(id uint64) error {
	return r.repository.Delete(id)
}

func (r *repository) ConvertResponse(p project.Project) ProjectResponse {
	var datasetCount int
	d, err := r.datasetRepo.GetByProject(p.ID)
	if err != nil {
		logger.Error("error getting dataset by project id")
	} else {
		datasetCount = len(d)
	}
	var pm = make([]uint64, 0)
	perms, totalPerms, err := r.repository.GetProjectPermissions(p.ID, project.Any, 0, 0)
	if err != nil {
		logger.Error("error getting project perm")
	}
	for i := range perms {
		if perms[i].Role == project.Manager {
			pm = append(pm, perms[i].UserID)
			break
		}
	}
	w, _ := r.workspaceRepo.GetByID(p.WorkspaceID)
	return ProjectResponse{
		ProjectBaseResponse: ProjectBaseResponse{
			ID:          p.ID,
			Title:       p.Title,
			Description: p.Description,
			Thumbnail:   p.Thumbnail,
			Color:       p.Color,
		},
		DatasetCount:   datasetCount,
		MemberCount:    totalPerms,
		Workspace:      w,
		ProjectManager: pm,
	}
}

func slice(projects []project.Project, offset, limit int) (filtered []project.Project) {
	l := len(projects)
	if offset >= l {
		return
	}
	if boundary := offset + limit; l >= boundary {
		return projects[offset:boundary]
	}
	return projects[offset:l]

}
