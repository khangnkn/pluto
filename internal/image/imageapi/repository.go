package imageapi

import (
	"bytes"
	"fmt"
	gimage "image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"net/url"

	"github.com/spf13/viper"

	"github.com/nkhang/pluto/internal/dataset"
	"github.com/nkhang/pluto/internal/image"
	"github.com/nkhang/pluto/pkg/errors"
	"github.com/nkhang/pluto/pkg/logger"
	"github.com/nkhang/pluto/pkg/objectstorage"
)

const (
	MaxLimitImage = 10
)

type Repository interface {
	GetByDatasetID(dID uint64, offset, limit int) ([]ImageResponse, error)
	UploadRequest(dID uint64, file *multipart.FileHeader) error
}

type repository struct {
	repo        image.Repository
	datasetRepo dataset.Repository
	storage     objectstorage.ObjectStorage
	conf        Config
}

func NewRepository(r image.Repository, s objectstorage.ObjectStorage, d dataset.Repository) *repository {
	var conf = Config{
		Scheme:   viper.GetString("minio.scheme"),
		Endpoint: viper.GetString("minio.endpoint"),
	}
	return &repository{
		repo:        r,
		storage:     s,
		datasetRepo: d,
		conf:        conf,
	}
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
	defer func() {
		err := file.Close()
		if err != nil {
			logger.Error(err)
		}
	}()

	var buf bytes.Buffer
	reader := io.TeeReader(file, &buf)

	img, _, err := gimage.Decode(reader)
	if err != nil {
		logger.Error("error decode image", err)
		return nil
	}

	collection := fmt.Sprintf("pluto-bucket-%d", dID)
	n, err := r.storage.PutImage(collection, header.Filename, &buf, header.Size)
	if err != nil {
		logger.Error("error putting to object storage", err)
		return err
	}
	logger.Infof("put image to object storage with %d bytes", n)

	w := img.Bounds().Max.X
	h := img.Bounds().Max.Y
	size := header.Size
	title := header.Filename
	u := r.getImageURL(collection, title)
	_, err = r.repo.CreateImage(title, u, w, h, size, dID)
	go func() {
		err := r.repo.InvalidateDatasetImage(dID)
		if err != nil {
			logger.Error("cannot invalidate dataset images", err)
		}
	}()
	return err
}

func (r *repository) getImageURL(collection, title string) string {
	return fmt.Sprintf("%s://%s/%s/%s", r.conf.Scheme, r.conf.Endpoint, collection, url.PathEscape(title))
}
