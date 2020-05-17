package imageapi

import (
	"github.com/nkhang/pluto/internal/image"
	"github.com/nkhang/pluto/pkg/errors"
)

const (
	MaxLimitImage = 10
)

type Repository interface {
	GetByDatasetID(dID uint64, offset, limit int) ([]ImageResponse, error)
}

type repository struct {
	repo image.Repository
}

func NewRepository(r image.Repository) *repository {
	return &repository{repo: r}
}

func (r *repository) GetByDatasetID(dID uint64, offset, limit int) ([]ImageResponse, error) {
	if limit > MaxLimitImage {
		return nil, errors.ImageTooManyRequest.NewWithMessage("too many image to query")
	}
	images, err := r.repo.GetByDataset(dID, offset, limit)
	if err != nil {
		return nil, err
	}
	responses := make([]ImageResponse, len(images))
	for i := range images {
		responses[i] = ToImageResponse(images[i])
	}
	return responses, nil
}
