package projectapi

import (
	"github.com/nkhang/pluto/internal/workspace/workspaceapi"
)

const (
	SrcAllProject = iota + 1
	SrcMyProject
	SrcOtherProject
)

type GetProjectRequest struct {
	Page     int `form:"page" binding:"required"`
	PageSize int `form:"page_size" binding:"required"`
	Source   int `form:"src" binding:"required"`
}

type CreateProjectRequest struct {
	Title       string `form:"title" json:"title"`
	Description string `form:"description" json:"description"`
	Color       string `form:"color" json:"color"`
}

type ProjectResponse struct {
	ProjectBaseResponse
	DatasetCount    int                                  `json:"dataset_count"`
	MemberCount     int                                  `json:"member_count"`
	ProjectManagers []uint64                             `json:"project_managers"`
	Workspace       workspaceapi.WorkspaceDetailResponse `json:"workspace"`
}

type ProjectBaseResponse struct {
	ID          uint64 `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Thumbnail   string `json:"thumbnail"`
	Color       string `json:"color"`
}

type GetProjectResponse struct {
	Total    int               `json:"total"`
	Projects []ProjectResponse `json:"projects"`
}

type UpdateProjectRequest struct {
	Title       string `form:"title" json:"title,omitempty"`
	Description string `form:"description" json:"description,omitempty"`
}
