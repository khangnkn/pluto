package labelapi

import (
	"github.com/nkhang/pluto/internal/label"
)

type LabelRequest struct {
	ProjectID uint64 `form:"project_id" binding:"required"`
}

type CreateLabelRequest struct {
	ProjectID uint64 `form:"project_id" binding:"required"`
	Name      string `form:"name" binding:"required"`
	Color     string `form:"color" binding:"required"`
	ToolID    uint64 `form:"tool_id" binding:"required"`
}

type LabelResponse struct {
	ID    uint64 `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
	Tool  string `json:"tool"`
}

func ToLabelResponse(l label.Label) LabelResponse {
	return LabelResponse{
		ID:    l.ID,
		Name:  l.Name,
		Color: l.Color,
		Tool:  l.Tool.Name,
	}
}
