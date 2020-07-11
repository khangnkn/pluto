package labelapi

import (
	"github.com/nkhang/pluto/internal/label"
	"github.com/nkhang/pluto/pkg/errors"
)

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

func (r *repository) CreateLabel(request CreateLabelRequest) error {
	errs := make([]error, 0)
	for _, req := range request.Labels {
		err := r.repository.CreateLabel(req.Name, req.Color, request.ProjectID, req.ToolID)
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) != 0 {
		return errors.LabelCannotCreate.NewWithMessage("cannot create labels")
	}
	return nil
}
