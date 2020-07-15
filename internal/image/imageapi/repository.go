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
	"github.com/nkhang/pluto/pkg/logger"
	"github.com/nkhang/pluto/pkg/objectstorage"
)

type Repository interface {
	GetImage(request GetImageRequest) (ImageResponse, error)
	GetByDatasetID(dID uint64, offset, limit int) ([]ImageResponse, error)
	UploadRequest(dID uint64, file []*multipart.FileHeader) []error
}

type repository struct {
	repo        image.Repository
	datasetRepo dataset.Repository
	storage     objectstorage.ObjectStorage
	conf        Config
}

func NewRepository(r image.Repository, s objectstorage.ObjectStorage, d dataset.Repository) *repository {
	var conf = Config{
		Scheme:     viper.GetString("minio.scheme"),
		Endpoint:   viper.GetString("minio.endpoint"),
		BucketName: viper.GetString("minio.bucketname"),
		BasePath:   viper.GetString("minio.basepath"),
	}
	return &repository{
		repo:        r,
		storage:     s,
		datasetRepo: d,
		conf:        conf,
	}
}

func (r *repository) GetImage(request GetImageRequest) (ImageResponse, error) {
	img, err := r.repo.Get(request.ID)
	if err != nil {
		return ImageResponse{}, err
	}
	return ToImageResponse(img), nil
}

func (r *repository) GetByDatasetID(dID uint64, offset, limit int) ([]ImageResponse, error) {
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

func (r *repository) UploadRequest(dID uint64, headers []*multipart.FileHeader) []error {
	errs := make([]error, 0)
	d, err := r.datasetRepo.Get(dID)
	if err != nil {
		return append(errs, err)
	}
	for _, header := range headers {
		err := r.createImage(d, header)
		if err != nil {
			errs = append(errs, err)
		}
	}
	go func() {
		err := r.repo.InvalidateDatasetImage(dID)
		if err != nil {
			logger.Error("cannot invalidate dataset images", err)
		}
	}()
	return errs
}

func (r *repository) getImageURL(collection, title string) string {
	return fmt.Sprintf("%s://%s/%s/%s", r.conf.Scheme, r.conf.BasePath, collection, url.PathEscape(title))
}

func (r *repository) createImage(d dataset.Dataset, h *multipart.FileHeader) error {
	file, err := h.Open()
	if err != nil {
		return err
	}
	defer func() {
		err := file.Close()
		if err != nil {
			logger.Error("error closing file", err)
		}
	}()

	var buf bytes.Buffer
	reader := io.TeeReader(file, &buf)

	img, _, err := gimage.Decode(reader)
	if err != nil {
		logger.Error("error decode image", err)
		return nil
	}

	path := fmt.Sprintf("%s/%d/%s", d.Project.Dir, d.ID, h.Filename)
	n, err := r.storage.PutImage(r.conf.BucketName, path, &buf, h.Size)
	if err != nil {
		logger.Error("error putting to object storage", err)
		return err
	}
	logger.Infof("put image to object storage with %d bytes", n)

	width := img.Bounds().Max.X
	height := img.Bounds().Max.Y
	size := h.Size
	u := r.getImageURL(r.conf.BucketName, path)
	_, err = r.repo.CreateImage(h.Filename, u, width, height, size, d.ID)
	return err
}
