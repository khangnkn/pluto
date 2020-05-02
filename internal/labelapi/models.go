package labelapi

import (
	"github.com/nkhang/pluto/internal/label"
	"github.com/nkhang/pluto/internal/toolapi"
)

type LabelResponse struct {
	ID   uint64               `json:"id"`
	Name string               `json:"name"`
	Tool toolapi.ToolResponse `json:"tool"`
}

func ToLabelResponse(l label.Label) LabelResponse {
	return LabelResponse{
		ID:   l.ID,
		Name: l.Name,
		Tool: toolapi.ToToolResponse(l.Tool),
	}
}
