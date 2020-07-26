package datasetapi

import (
	"net/url"
	"strings"

	"github.com/nkhang/pluto/internal/project"

	"github.com/nkhang/pluto/internal/dataset"
	"github.com/nkhang/pluto/internal/image"
	"github.com/nkhang/pluto/internal/tool/enc"
	"github.com/nkhang/pluto/pkg/errors"
	"github.com/nkhang/pluto/pkg/logger"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

type Repository interface {
	GetByID(dID uint64) (DatasetResponse, error)
	GetByProjectID(pID uint64) ([]DatasetResponse, error)
	CreateDataset(title, description string, pID uint64) (DatasetResponse, error)
	Delete(id uint64) error
	CloneDataset(projectID uint64, token string) (DatasetResponse, error)
	GetLink(datasetID uint64) (string, error)
	ParseLink(link string) (GetLinkResponse, error)
}

type repository struct {
	repository  dataset.Repository
	imgRepo     image.Repository
	projectRepo project.Repository
	baseURL     string
	secret      []byte
}

func NewRepository(r dataset.Repository, imgRepo image.Repository, p project.Repository) *repository {
	secret := viper.GetString("getlink.secret")
	baseURL := viper.GetString("getlink.baseurl")
	if secret == "" || baseURL == "" {
		logger.Panic("secret empty")
	}
	return &repository{
		repository:  r,
		imgRepo:     imgRepo,
		baseURL:     baseURL,
		projectRepo: p,
		secret:      []byte(secret),
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

func (r *repository) CloneDataset(dest uint64, token string) (DatasetResponse, error) {
	token = strings.TrimPrefix(token, "/")
	idStr, err := enc.Decrypt(r.secret, token)
	if err != nil {
		return DatasetResponse{}, errors.DatasetLinkCannotParse.NewWithMessage("cannot parse provided link")
	}
	datasetID, err := cast.ToUint64E(idStr)
	if err != nil {
		return DatasetResponse{}, errors.DatasetLinkCannotParse.NewWithMessage("cannot parse dataset ID")
	}
	images, err := r.imgRepo.GetAllImageByDataset(datasetID)
	if err != nil {
		logger.Error("getting all image error", err)
		return DatasetResponse{}, nil
	}
	err = r.imgRepo.BulkInsert(images, dest)
	if err != nil {
		return DatasetResponse{}, err
	}
	if len(images) != 0 {
		cloned, err := r.repository.Update(dest, map[string]interface{}{
			"thumbnail": images[0].Thumbnail,
		})
		if err == nil {
			err = r.projectRepo.PickThumbnail(cloned.ProjectID)
			if err != nil {
				logger.Errorf("[DATASET-API] - error picking project thumbnail. err %v", err)
			}
			return r.ToDatasetResponse(cloned), nil
		}
	}
	cloned, err := r.repository.Get(dest)
	if err != nil {
		return DatasetResponse{}, err
	}
	return r.ToDatasetResponse(cloned), nil
}

func (r *repository) Delete(id uint64) error {
	d, err := r.repository.Get(id)
	if err != nil {
		return err
	}
	err = r.repository.DeleteDataset(id)
	if err != nil {
		return err
	}
	err = r.projectRepo.PickThumbnail(d.ProjectID)
	if err != nil {
		logger.Errorf("[DATASET-API] - error picking project thumbnail. err %v", err)
	}
	return nil
}

func (r *repository) GetLink(datasetID uint64) (string, error) {
	token, err := enc.Encrypt(r.secret, cast.ToString(datasetID))
	if err != nil {
		return "", errors.DatasetLinkCannotParse.NewWithMessage("cannot get token")
	}
	return r.baseURL + "/" + token, nil
}

func (r *repository) ParseLink(link string) (GetLinkResponse, error) {
	URL, err := url.Parse(link)
	if err != nil {
		return GetLinkResponse{}, errors.DatasetLinkCannotParse.NewWithMessage("cannot parse provided link")
	}
	token := URL.Path
	token = strings.TrimPrefix(token, "/")
	idStr, err := enc.Decrypt(r.secret, token)
	if err != nil {
		return GetLinkResponse{}, errors.DatasetLinkCannotParse.NewWithMessage("cannot parse provided link")
	}
	id, err := cast.ToUint64E(idStr)
	if err != nil {
		return GetLinkResponse{}, errors.DatasetLinkCannotParse.NewWithMessage("cannot parse provided link")
	}
	dataset, err := r.repository.Get(id)
	if err != nil {
		return GetLinkResponse{}, err
	}
	return r.ToDatasetResponse(dataset).WithToken(token), nil
}
