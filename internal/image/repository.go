package image

import (
	"github.com/nkhang/pluto/internal/rediskey"
	"github.com/nkhang/pluto/pkg/cache"
	"github.com/nkhang/pluto/pkg/errors"
	"github.com/nkhang/pluto/pkg/logger"
)

type Repository interface {
	Get(id uint64) (Image, error)
	GetByDataset(dID uint64, offset, limit int) (imgs []Image, err error)
	CreateImage(name string, w, h int, dataset_id uint64) (Image, error)
}

type repository struct {
	dbRepo    DBRepository
	cacheRepo cache.Cache
}

func NewRepository(d DBRepository, c cache.Cache) *repository {
	return &repository{
		dbRepo:    d,
		cacheRepo: c,
	}
}

func (r *repository) Get(id uint64) (img Image, err error) {
	key := rediskey.ImageByID(id)
	err = r.cacheRepo.Get(key, &img)
	if err == nil {
		logger.Infof("cache hit for image %d", id)
		return
	}
	if errors.Type(err) == errors.CacheNotFound {
		logger.Infof("cache miss for image %d", id)
	} else {
		logger.Errorf("error getting image %d from cache", id)
	}
	img, err = r.dbRepo.Get(id)
	if err != nil {
		return
	}
	go func() {
		err := r.cacheRepo.Set(key, img)
		if err != nil {
			logger.Error(err)
		}
	}()
	return
}

func (r *repository) GetByDataset(dID uint64, offset, limit int) (images []Image, err error) {
	key := rediskey.ImageByDatasetID(dID, offset, limit)
	err = r.cacheRepo.Get(key, &images)
	if err == nil {
		logger.Infof("cache hit for images by dataset %d", dID)
		return
	}
	if errors.Type(err) == errors.CacheNotFound {
		logger.Infof("cache miss for images by dataset %d", dID)
	} else {
		logger.Errorf("error getting images by dataset %d from cache", dID)
	}
	images, err = r.dbRepo.GetByDataset(dID, offset, limit)
	if err != nil {
		return
	}
	go func() {
		err := r.cacheRepo.Set(key, images)
		if err != nil {
			logger.Error(err)
		}
	}()
	return
}

func (r *repository) CreateImage(name string, w, h int, dataset_id uint64) (Image, error) {
	return r.dbRepo.CreateImage(name, w, h, dataset_id)
}
