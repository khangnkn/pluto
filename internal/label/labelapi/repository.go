package labelapi

import "github.com/nkhang/pluto/internal/label"

type Repository interface {
	GetByProject(pID uint64) ([]LabelResponse, error)
	CreateLabel(r CreateLabelRequest) error
}

type repository struct {
	repository label.Repository
}

func NewRepository(r label.Repository) *repository {
	return &repository{
		repository: r,
	}
}

func (r *repository) GetByProject(pID uint64) ([]LabelResponse, error) {
	labels, err := r.repository.GetByProjectId(pID)
	if err != nil {
		return nil, err
	}
	responses := make([]LabelResponse, len(labels))
	for i := range labels {
		responses[i] = ToLabelResponse(labels[i])
	}
	return responses, nil
}

func (r *repository) CreateLabel(req CreateLabelRequest) error {
	return r.repository.CreateLabel(req.Name, req.Color, req.ProjectID, req.ToolID)
}
