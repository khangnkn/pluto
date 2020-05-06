package projectapi

import (
	"github.com/nkhang/pluto/internal/project"
	"github.com/nkhang/pluto/pkg/logger"
)

type Repository interface {
	GetByID(wID uint64, pID uint64) (ProjectResponse, error)
	GetByWorkspaceID(id uint64) ([]ProjectResponse, error)
}

type repository struct {
	repository project.Repository
}

func NewRepository(r project.Repository) *repository {
	return &repository{repository: r}
}

func (r *repository) GetByID(wID uint64, pID uint64) (ProjectResponse, error) {
	p, err := r.repository.Get(wID, pID)
	if err != nil {
		return ProjectResponse{}, err
	}
	return ProjectResponse{
		ID:          p.ID,
		Title:       p.Title,
		Description: p.Description,
	}, nil
}

func (r *repository) GetByWorkspaceID(id uint64) ([]ProjectResponse, error) {
	projects, err := r.repository.GetByWorkspaceID(id)
	if err != nil {
		return nil, err
	}
	logger.Info(projects)
	responses := make([]ProjectResponse, len(projects))
	for i := range projects {
		responses[i] = toProjectResponse(projects[i])
	}
	return responses, nil
}
