package labelapi

import (
	"github.com/nkhang/pluto/internal/label"
	"github.com/nkhang/pluto/internal/tool/toolapi"
)

type LabelRequest struct {
	ProjectID uint64 `json:"project_id" form:"project_id"`
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
