package toolapi

import (
	"github.com/nkhang/pluto/internal/tool"
)

type ToolRepository interface {
	GetAll() ([]tool.Tool, error)
}

type Repository ToolRepository
type repository struct {
	toolRepo Repository
}

func NewRepository(r ToolRepository) *repository {
	return &repository{
		toolRepo: r,
	}
}

func (r *repository) GetAll() ([]tool.Tool, error) {
	return r.toolRepo.GetAll()
}
