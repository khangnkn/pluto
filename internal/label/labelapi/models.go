package labelapi

import (
	"github.com/nkhang/pluto/internal/label"
	"github.com/nkhang/pluto/internal/tool/toolapi"
)

type LabelRequest struct {
	ProjectID uint64 `form:"project_id" binding:"required"`
}

type CreateLabelRequest struct {
	ProjectID uint64              `json:"project_id" form:"project_id" binding:"required"`
	Labels    []CreateLabelObject `json:"labels" form:"labels"`
}

type CreateLabelObject struct {
	Name   string `json:"name" form:"name" binding:"required"`
	Color  string `json:"color" form:"color" binding:"required"`
	ToolID uint64 `json:"tool" form:"tool" binding:"required"`
}

type LabelResponse struct {
	ID    uint64               `json:"id"`
	Name  string               `json:"name"`
	Color string               `json:"color"`
	Tool  toolapi.ToolResponse `json:"tool"`
}

func ToLabelResponse(l label.Label) LabelResponse {
	return LabelResponse{
		ID:    l.ID,
		Name:  l.Name,
		Color: l.Color,
		Tool:  toolapi.ToToolResponse(l.Tool),
	}
}
