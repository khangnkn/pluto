package permissionapi

import "github.com/nkhang/pluto/internal/project"

type CreatePermRequest struct {
	Members []CreatePermObject `json:"members"`
}

type CreatePermObject struct {
	UserID uint64       `form:"user_id" json:"user_id" binding:"required"`
	Role   project.Role `form:"role" json:"role" binding:"required"`
}

type PermissionResponse struct {
	Total   int                `json:"total"`
	Members []PermissionObject `json:"members"`
}

type PermissionObject struct {
	CreatedAt int64        `json:"created_at"`
	UserID    uint64       `json:"user_id"`
	Role      project.Role `json:"role"`
}

type UpdatePermissionRequest struct {
	UserID uint64       `form:"user_id" json:"user_id" binding:"required"`
	Role   project.Role `form:"role" json:"role" binding:"required"`
}
