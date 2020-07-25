package permissionapi

import (
	"github.com/nkhang/pluto/internal/project"
	"github.com/nkhang/pluto/internal/project/projectapi"
	"github.com/nkhang/pluto/pkg/annotation"
	"github.com/nkhang/pluto/pkg/errors"
	"github.com/nkhang/pluto/pkg/logger"
	"github.com/nkhang/pluto/pkg/util/clock"
)

type Repository interface {
	Create(projectID uint64, req CreatePermRequest) (projectapi.ProjectResponse, error)
	GetList(projectID uint64) (PermissionResponse, error)
	Update(projectID uint64, req UpdatePermissionRequest) (PermissionObject, error)
	Delete(projectID, userID uint64) error
}

type repository struct {
	repository        project.Repository
	projectRepo       projectapi.Repository
	annotationService annotation.Service
}

func NewProjectPermissionAPIRepository(r project.Repository, p projectapi.Repository, ann annotation.Service) *repository {
	return &repository{
		repository:        r,
		projectRepo:       p,
		annotationService: ann,
	}
}

func (r *repository) Create(projectID uint64, req CreatePermRequest) (resp projectapi.ProjectResponse, err error) {
	var errs = make([]error, 0)
	for _, p := range req.Members {
		_, err := r.repository.GetPermission(p.UserID, projectID)
		if err == nil {
			continue
		}
		if p.Role == project.Admin { //role Admin
			continue
		}
		_, err = r.repository.CreatePermission(projectID, p.UserID, p.Role)
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) != 0 {
		err = errors.ProjectPermissionCreatingError.NewWithMessageF("error creating %d permissions", len(errs))
		return
	}
	err = r.annotationService.UpdateProject(projectID)
	if err != nil {
		logger.Errorf("from project: push project to annotation server failed %v", err)
	}
	resp, err = r.projectRepo.GetByID(projectID)
	return
}

func (r *repository) GetList(projectID uint64) (PermissionResponse, error) {
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

func (r *repository) Update(projectID uint64, req UpdatePermissionRequest) (PermissionObject, error) {
	if req.Role != project.Member && req.Role != project.Manager {
		return PermissionObject{}, errors.ProjectPermissionCannotUpdate.NewWithMessage("role not supported")
	}
	perm, err := r.repository.UpdatePermission(projectID, req.UserID, req.Role)
	if err != nil {
		return PermissionObject{}, err
	}
	return convertPermissionObject(perm), nil
}

func convertPermissionObject(perm project.Permission) PermissionObject {
	return PermissionObject{
		CreatedAt: clock.UnixMillisecondFromTime(perm.CreatedAt),
		UserID:    perm.UserID,
		Role:      perm.Role,
	}
}

func (r *repository) Delete(projectID, userID uint64) error {
	return r.repository.DeletePermission(userID, projectID)
}
