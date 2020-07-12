package datasetapi

import (
	"github.com/nkhang/pluto/internal/dataset"
	"github.com/nkhang/pluto/internal/image"
	"github.com/nkhang/pluto/pkg/logger"
)

type Repository interface {
	GetByID(dID uint64) (DatasetResponse, error)
	GetByProjectID(pID uint64) ([]DatasetResponse, error)
	CreateDataset(title, description string, pID uint64) (DatasetResponse, error)
	Delete(id uint64) error
	CloneDataset(projectID uint64, datasetID uint64) (DatasetResponse, error)
}

type repository struct {
	repository dataset.Repository
	imgRepo    image.Repository
}

func NewRepository(r dataset.Repository, imgRepo image.Repository) *repository {
	return &repository{
		repository: r,
		imgRepo:    imgRepo,
	}
}

func (r *repository) GetByID(dID uint64) (DatasetResponse, error) {
	d, err := r.repository.Get(dID)
	if err != nil {
		return DatasetResponse{}, err
	}
	return r.ToDatasetResponse(d), nil
}

func (r *repository) GetByProjectID(pID uint64) ([]DatasetResponse, error) {
	datasets, err := r.repository.GetByProject(pID)
	if err != nil {
		return nil, err
	}
	responses := make([]DatasetResponse, len(datasets))
	for i := range datasets {
		responses[i] = r.ToDatasetResponse(datasets[i])
	}
	return responses, nil
}

func (r *repository) CreateDataset(title, description string, pID uint64) (DatasetResponse, error) {
	d, err := r.repository.CreateDataset(title, description, pID)
	if err != nil {
		return DatasetResponse{}, err
	}
	return r.ToDatasetResponse(d), nil
}

func (r *repository) CloneDataset(projectID uint64, datasetID uint64) (DatasetResponse, error) {
	origin, err := r.repository.Get(datasetID)
	if err != nil {
		logger.Errorf("error getting dataset %d, error %v", datasetID, err)
		return DatasetResponse{}, err
	}
	images, err := r.imgRepo.GetAllImageByDataset(datasetID)
	if err != nil {
		logger.Error("getting all image error", err)
		return DatasetResponse{}, nil
	}
	cloned, err := r.repository.CreateDataset(origin.Title, origin.Description, projectID)
	if err != nil {
		logger.Errorf("cannot creating dataset")
		return DatasetResponse{}, err
	}
	logger.Info("clone dataset successfully", cloned)
	err = r.imgRepo.BulkInsert(images, cloned.ID)
	if err != nil {
		logger.Errorf("error inserting images for dataset %d. now rollback creating", cloned.ID)
		go func() {
			err := r.repository.DeleteDataset(cloned.ID)
			if err != nil {
				logger.Errorf("cannot delete uncompleted dataset %d, error", cloned.ID, err)
			}
		}()
		return DatasetResponse{}, err
	}
	return r.ToDatasetResponse(cloned), nil
}

func (r *repository) Delete(id uint64) error {
	return r.repository.DeleteDataset(id)
}
