package toolapi

import "github.com/nkhang/pluto/internal/tool"

type ToolRepository interface {
	GetAll() ([]tool.Tool, error)
}
type repository struct {
	toolRepo ToolRepository
}

func NewRepository() *repository {
	return &repository{}
}

func (r *repository) GetAll() ([]tool.Tool, error) {

}
