package workspaceapi

import (
	"github.com/nkhang/pluto/internal/workspace"
	"github.com/nkhang/pluto/pkg/util/clock"
)

type GetByUserIDRequest struct {
	Page     int `form:"page" binding:"required"`
	PageSize int `form:"page_size" binding:"required"`
	Source   int `form:"src" binding:"required"`
}

type CreateWorkspaceRequest struct {
	Title       string   `form:"title" json:"title"`
	Description string   `form:"description" json:"description"`
	Color       string   `form:"color" json:"color"`
	Members     []uint64 `form:"members" json:"members"`
}

type UpdateWorkspaceRequest struct {
	Title       string `form:"title" json:"title,omitempty"`
	Description string `form:"description" json:"description,omitempty"`
}

type GetByUserResponse struct {
	Total      int                       `json:"total"`
	Workspaces []WorkspaceDetailResponse `json:"workspaces"`
}

type WorkspaceDetailResponse struct {
	WorkspaceResponse
	ProjectCount int    `json:"project_count"`
	MemberCount  int    `json:"member_count"`
	Admin        uint64 `json:"admin"`
}

type WorkspaceResponse struct {
	ID          uint64 `json:"id"`
	Updated     int64  `json:"updated"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

func ToWorkspaceInfoResponse(w workspace.Workspace) WorkspaceResponse {
	return WorkspaceResponse{
		ID:          w.ID,
		Updated:     clock.UnixMillisecondFromTime(w.UpdatedAt),
		Title:       w.Title,
		Description: w.Description,
		Color:       w.Color,
	}
}
