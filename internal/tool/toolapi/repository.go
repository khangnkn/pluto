package toolapi

import "github.com/nkhang/pluto/internal/tool"

type Repository interface {
	GetAll() ([]ToolResponse, error)
}

type repository struct {
	toolRepo tool.Repository
}

func NewRepository(r tool.Repository) *repository {
	return &repository{
		toolRepo: r,
	}
}

func (r *repository) GetAll() ([]ToolResponse, error) {
	tools, err := r.toolRepo.GetAll()
	if err != nil {
		return nil, err
	}
	responses := make([]ToolResponse, len(tools))
	for i := range tools {
		responses[i] = ToToolResponse(tools[i])
	}
	return responses, nil
}
