package projectapi

import "github.com/nkhang/pluto/internal/project"

type Repository interface {
	GetByID(id uint64) (ProjectResponse, error)
}

type repository struct {
	repository project.Repository
}

func NewRepository(r project.Repository) *repository {
	return &repository{repository: r}
}

func (r *repository) GetByID(id uint64) (ProjectResponse, error) {
	p, err := r.repository.Get(id)
	if err != nil {
		return ProjectResponse{}, err
	}
	return ProjectResponse{
		ID:          p.ID,
		Title:       p.Title,
		Description: p.Description,
	}, nil
}
