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

	"github.com/nkhang/pluto/internal/dataset/datasetapi"
	"github.com/nkhang/pluto/pkg/util/clock"

	"github.com/nkhang/pluto/internal/project"

	"github.com/spf13/viper"

	"github.com/nkhang/pluto/internal/dataset"
	"github.com/nkhang/pluto/internal/image"
	"github.com/nkhang/pluto/pkg/logger"
	"github.com/nkhang/pluto/pkg/objectstorage"
)

type Repository interface {
	GetImage(request GetImageRequest) (ImageResponse, error)
	GetByDatasetID(dID uint64, offset, limit int) ([]ImageResponse, error)
	UploadRequest(dID uint64, headers []*multipart.FileHeader) (datasetapi.DatasetResponse, []error)
}

type repository struct {
	repo        image.Repository
	datasetRepo dataset.Repository
	projectRepo project.Repository
	storage     objectstorage.ObjectStorage
	conf        Config
}

func NewRepository(r image.Repository, s objectstorage.ObjectStorage, d dataset.Repository, p project.Repository) *repository {
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
		projectRepo: p,
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

func (r *repository) UploadRequest(dID uint64, headers []*multipart.FileHeader) (datasetapi.DatasetResponse, []error) {
	errs := make([]error, 0)
	d, err := r.datasetRepo.Get(dID)
	if err != nil {
		return datasetapi.DatasetResponse{}, append(errs, err)
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
		if len(headers) > len(errs) {
			img, err := r.repo.GetAllImageByDataset(dID)
			if err != nil && len(img) > 0 {
				logger.Error("cannot get image to set to dataset")
				return
			}
			d, err := r.datasetRepo.Update(dID, map[string]interface{}{
				"thumbnail": img[0].URL,
			})
			if err != nil {
				logger.Errorf("cannot update dataset %d thumbnail", d.ID)
				return
			}
			_, err = r.projectRepo.UpdateProject(d.ProjectID, map[string]interface{}{
				"thumbnail": img[0].URL,
			})
			if err != nil {
				logger.Errorf("cannot update project %d thumbnail", d.ID)
				return
			}
		}
	}()
	d, err = r.datasetRepo.Get(dID)
	if err != nil {
		logger.Errorf("cannot get dataset after upload task %d", dID)
		return datasetapi.DatasetResponse{}, errs
	}
	imgs, err := r.repo.GetAllImageByDataset(d.ID)
	if err != nil {
		return datasetapi.DatasetResponse{}, append(errs, err)
	}
	return datasetapi.DatasetResponse{
		ID:          d.ID,
		Title:       d.Title,
		Description: d.Description,
		Thumbnail:   d.Thumbnail,
		ProjectID:   d.ProjectID,
		ImageCount:  len(imgs),
		UpdatedAt:   clock.UnixMillisecondFromTime(d.UpdatedAt),
	}, errs
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
	prj, err := r.projectRepo.Get(d.ProjectID)
	if err != nil {
		return err
	}
	path := fmt.Sprintf("%s/%d/%s", prj.Dir, d.ID, h.Filename)
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
