package imageapi

import (
	gimage "image"
	_ "image/png"
	_ "image/jpeg"
	"mime/multipart"

	"github.com/nkhang/pluto/internal/image"
	"github.com/nkhang/pluto/pkg/errors"
	"github.com/nkhang/pluto/pkg/logger"
)

const (
	MaxLimitImage = 10
)

type Repository interface {
	GetByDatasetID(dID uint64, offset, limit int) ([]ImageResponse, error)
	UploadRequest(dID uint64, file *multipart.FileHeader) error
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

func (r *repository) UploadRequest(dID uint64, header *multipart.FileHeader) error {
	file, err := header.Open()
	if err != nil {
		return err
	}
	defer file.Close()
	image, t, err := gimage.Decode(file)
	if err != nil {
		logger.Error(err)
		return err
	}
	logger.Info(t)
	logger.Info(image.Bounds().Dx, image.Bounds().Dy)
	return nil
}
