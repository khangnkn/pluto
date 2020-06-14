package workspaceapi

import (
	"github.com/nkhang/pluto/internal/workspace"
)

type GetByUserIDRequest struct {
	UserID uint64 `form:"user_id" binding:"required"`
}

type CreateWorkspaceRequest struct {
	UserID      uint64 `form:"user_id" json:"user_id" binding:"required"`
	Title       string `form:"title" json:"title"`
	Description string `form:"title" json:"description"`
}

type WorkspaceInfoResponse struct {
	ID          uint64 `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func toWorkspaceInfoResponse(w workspace.Workspace) WorkspaceInfoResponse {
	return WorkspaceInfoResponse{
		ID:          w.ID,
		Title:       w.Title,
		Description: w.Description,
	}
}
