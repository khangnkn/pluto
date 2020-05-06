package datasetapi

import "github.com/nkhang/pluto/internal/dataset"

type Repository interface {
	GetByID(dID uint64) (DatasetResponse, error)
	GetByProjectID(pID uint64) ([]DatasetResponse, error)
}

type repository struct {
	repository dataset.Repository
}

func NewRepository(r dataset.Repository) *repository {
	return &repository{
		repository: r,
	}
}

func (r *repository) GetByID(dID uint64) (DatasetResponse, error) {
	d, err := r.repository.Get(dID)
	if err != nil {
		return DatasetResponse{}, err
	}
	return ToDatasetResponse(d), nil
}

func (r *repository) GetByProjectID(pID uint64) ([]DatasetResponse, error) {
	datasets, err := r.repository.GetByProject(pID)
	if err != nil {
		return nil, err
	}
	responses := make([]DatasetResponse, len(datasets))
	for i := range datasets {
		responses[i] = ToDatasetResponse(datasets[i])
	}
	return responses, nil
}
