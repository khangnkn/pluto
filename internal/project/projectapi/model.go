package projectapi

import "github.com/nkhang/pluto/internal/project"

type GetProjectParam struct {
	WorkspaceID uint64 `form:"workspace_id"`
}
type ProjectResponse struct {
	ID           uint64 `json:"id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	DatasetCount int    `json:"dataset_count"`
	MemberCount  int    `json:"member_count"`
}

func toProjectResponse(p project.Project) ProjectResponse {
	return ProjectResponse{
		ID:          p.ID,
		Title:       p.Title,
		Description: p.Description,
	}
}
