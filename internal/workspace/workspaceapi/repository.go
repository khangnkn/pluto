package workspaceapi

import (
	"encoding/json"

	"github.com/nkhang/pluto/internal/project"
	"github.com/nkhang/pluto/internal/workspace"
	"github.com/nkhang/pluto/pkg/errors"
	"github.com/nkhang/pluto/pkg/logger"
	"github.com/nkhang/pluto/pkg/util/paging"
)

type Repository interface {
	GetByID(id uint64) (WorkspaceResponse, error)
	GetByUserID(request GetByUserIDRequest) (GetByUserResponse, error)
	CreateWorkspace(p CreateWorkspaceRequest) (WorkspaceResponse, error)
	UpdateWorkspace(id uint64, request UpdateWorkspaceRequest) (WorkspaceResponse, error)
}

type repository struct {
	workspaceRepository workspace.Repository
	projectRepo         project.Repository
}

func NewRepository(workspaceRepo workspace.Repository, projectRepo project.Repository) *repository {
	return &repository{
		workspaceRepository: workspaceRepo,
		projectRepo:         projectRepo,
	}
}

func (r *repository) GetByID(id uint64) (WorkspaceResponse, error) {
	w, err := r.workspaceRepository.Get(id)
	if err != nil {
		return WorkspaceResponse{}, err
	}
	return r.convertResponse(
		w), nil
}

func (r *repository) GetByUserID(request GetByUserIDRequest) (GetByUserResponse, error) {
	offset, limit := paging.Parse(request.Page, request.PageSize)
	var (
		workspaces []workspace.Workspace
		total      int
		err        error
	)
	switch request.Source {
	case 1:
		workspaces, total, err = r.workspaceRepository.GetByUserID(request.UserID, workspace.Any, offset, limit)
		if err != nil {
			return GetByUserResponse{}, err
		}
	case 2:
		workspaces, total, err = r.workspaceRepository.GetByUserID(request.UserID, workspace.Admin, offset, limit)
		if err != nil {
			return GetByUserResponse{}, err
		}
	case 3:
		workspaces, total, err = r.workspaceRepository.GetByUserID(request.UserID, workspace.Member, offset, limit)
		if err != nil {
			return GetByUserResponse{}, err
		}
	default:
		return GetByUserResponse{}, errors.BadRequest.NewWithMessage("unsupported src")
	}
	responses := make([]WorkspaceResponse, len(workspaces))
	for i := range workspaces {
		responses[i] = r.convertResponse(workspaces[i])
	}
	return GetByUserResponse{
		Total:      total,
		Workspaces: responses,
	}, nil
}

func (r *repository) CreateWorkspace(p CreateWorkspaceRequest) (WorkspaceResponse, error) {
	w, err := r.workspaceRepository.Create(p.UserID, p.Title, p.Description)
	if err != nil {
		return WorkspaceResponse{}, err
	}
	go r.workspaceRepository.InvalidateForUser(p.UserID)
	response := r.convertResponse(w)
	return response, nil

}

func (r *repository) convertResponse(w workspace.Workspace) WorkspaceResponse {
	logger.Infof("%+v", w)
	_, projectCount, err := r.projectRepo.GetByWorkspaceID(w.ID, 0, 0)
	if err != nil {
		logger.Error("cannot get all projects by workspace")
		projectCount = 0
	}
	_ = projectCount
	_, permissionCount, err := r.workspaceRepository.GetPermission(w.ID, workspace.Any, 0, 0)
	if err != nil {
		logger.Error("cannot get all permissions by workspace")
		permissionCount = 0
	}
	var admin uint64
	perm, _, err := r.workspaceRepository.GetPermission(w.ID, workspace.Admin, 0, 0)
	if err == nil && len(perm) != 0 {
		admin = perm[0].UserID
	} else {
		logger.Error("get admin error")
	}
	return WorkspaceResponse{
		ID:           w.ID,
		Title:        w.Title,
		Description:  w.Description,
		ProjectCount: projectCount,
		MemberCount:  permissionCount,
		Admin:        admin,
	}
}

func (r *repository) UpdateWorkspace(id uint64, request UpdateWorkspaceRequest) (WorkspaceResponse, error) {
	var changes = make(map[string]interface{})
	b, _ := json.Marshal(&request)
	_ = json.Unmarshal(b, &changes)
	logger.Info(changes)
	w, err := r.workspaceRepository.UpdateWorkspace(id, changes)
	if err != nil {
		return WorkspaceResponse{}, nil
	}
	return r.convertResponse(w), nil
}
