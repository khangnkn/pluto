package projectapi

import "github.com/nkhang/pluto/internal/project"

type ProjectResponse struct {
	ID          uint64 `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func toProjectResponse(p project.Project) ProjectResponse {
	return ProjectResponse{
		ID:          p.ID,
		Title:       p.Title,
		Description: p.Description,
	}
}
