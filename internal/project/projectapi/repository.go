package projectapi

import (
	"github.com/nkhang/pluto/internal/dataset"
	"github.com/nkhang/pluto/internal/project"
	"github.com/nkhang/pluto/pkg/logger"
)

type Repository interface {
	GetByID(pID uint64) (ProjectResponse, error)
	GetByWorkspaceID(id uint64) ([]ProjectResponse, error)
}

type repository struct {
	repository  project.Repository
	datasetRepo dataset.Repository
}

func NewRepository(r project.Repository, dr dataset.Repository) *repository {
	return &repository{
		repository:  r,
		datasetRepo: dr,
	}
}

func (r *repository) GetByID(pID uint64) (ProjectResponse, error) {
	p, err := r.repository.Get(pID)
	if err != nil {
		return ProjectResponse{}, err
	}
	return r.convertResponse(p), nil
}

func (r *repository) GetByWorkspaceID(id uint64) ([]ProjectResponse, error) {
	projects, err := r.repository.GetByWorkspaceID(id)
	if err != nil {
		return nil, err
	}
	logger.Info(projects)
	responses := make([]ProjectResponse, len(projects))
	for i := range projects {
		responses[i] = r.convertResponse(projects[i])
	}
	return responses, nil
}

func (r *repository) convertResponse(p project.Project) ProjectResponse {
	var datasetCount int
	var memberCount int
	d, err := r.datasetRepo.GetByProject(p.ID)
	if err != nil {
		logger.Error("error getting project by project id")
	} else {
		datasetCount = len(d)
	}
	m, err := r.repository.GetProjectPermission(p.ID)
	if err != nil {
		logger.Error("error getting project perm")
	} else {
		memberCount = len(m)
	}
	return ProjectResponse{
		ID:           p.ID,
		Title:        p.Title,
		Description:  p.Description,
		DatasetCount: datasetCount,
		MemberCount:  memberCount,
	}
}
