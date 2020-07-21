package permissionapi

import (
	"github.com/nkhang/pluto/internal/project"
	"github.com/nkhang/pluto/pkg/errors"
	"github.com/nkhang/pluto/pkg/util/clock"
)

type Repository interface {
	Create(projectID uint64, req CreatePermRequest) error
	GetList(projectID uint64) (PermissionResponse, error)
	Update(projectID uint64, req UpdatePermissionRequest) (PermissionObject, error)
	Delete(projectID, userID uint64) error
}

type repository struct {
	repository project.Repository
}

func NewProjectPermissionAPIRepository(r project.Repository) *repository {
	return &repository{
		repository: r,
	}
}

func (r *repository) Create(projectID uint64, req CreatePermRequest) error {
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
		return errors.ProjectPermissionCreatingError.NewWithMessageF("error creating %d permissions", len(errs))
	}
	return nil
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
	var role project.Role
	switch req.Role {
	case 2:
		role = project.Manager
	case 3:
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
