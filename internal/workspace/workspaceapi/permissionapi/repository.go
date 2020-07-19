package permissionapi

import (
	"github.com/nkhang/pluto/internal/workspace"
	"github.com/nkhang/pluto/pkg/util/clock"
	"github.com/nkhang/pluto/pkg/util/paging"
)

type Repository interface {
	CreatePermissions(id uint64, request CreatePermsRequest) error
	GetPermissions(workspaceID uint64, request GetPermsRequest) (GetPermissionResponse, error)
	DeletePermission(workspaceID uint64, userID uint64) error
}

type repository struct {
	workspaceRepo workspace.Repository
}

func NewRepository(workspaceRepo workspace.Repository) *repository {
	return &repository{workspaceRepo: workspaceRepo}
}

func (r *repository) CreatePermissions(id uint64, request CreatePermsRequest) error {
	return r.workspaceRepo.CreatePermission(id, request.UserIDs, workspace.Member)
}

func (r *repository) GetPermissions(workspaceID uint64, request GetPermsRequest) (GetPermissionResponse, error) {
	offset, limit := paging.Parse(request.Page, request.PageSize)
	perms, count, err := r.workspaceRepo.GetPermission(workspaceID, workspace.Any, offset, limit)
	if err != nil {
		return GetPermissionResponse{}, err
	}
	responses := make([]PermissionResponse, len(perms))
	for i := range perms {
		responses[i] = r.ToPermissionResponse(perms[i])
	}
	return GetPermissionResponse{
		Total:   count,
		Members: responses,
	}, nil
}

func (r *repository) DeletePermission(workspaceID uint64, userID uint64) error {
	return r.workspaceRepo.DeletePermission(workspaceID, userID)
}

func (r *repository) ToPermissionResponse(perm workspace.Permission) PermissionResponse {
	return PermissionResponse{
		CreatedAt: clock.UnixMillisecondFromTime(perm.CreatedAt),
		UserID:    perm.UserID,
		Role:      int32(perm.Role),
	}
}
