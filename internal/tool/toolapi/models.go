package toolapi

import "github.com/nkhang/pluto/internal/tool"

type ToolResponse struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

func ToToolResponse(t tool.Tool) ToolResponse {
	return ToolResponse{
		ID:   t.ID,
		Name: t.Name,
	}
}
