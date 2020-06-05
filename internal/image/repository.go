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
	GetAllImageByDataset(dID uint64) ([]Image, error)
	CreateImage(title, url string, w, h int, size int64, dataset_id uint64) (Image, error)
	BulkInsert(images []Image, dID uint64) error
	InvalidateDatasetImage(dID uint64) error
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

func (r *repository) CreateImage(title, url string, w, h int, size int64, dataset_id uint64) (Image, error) {
	return r.dbRepo.CreateImage(title, url, w, h, size, dataset_id)
}

func (r *repository) InvalidateDatasetImage(dID uint64) error {
	pattern := rediskey.ImageByDatasetIDAllKeys(dID)
	keys, err := r.cacheRepo.Keys(pattern)
	if err != nil {
		logger.Error("error getting all keys from redis", err)
		return err
	}
	logger.Infof("the following keys will be deleted: %v", keys)
	return r.cacheRepo.Del(keys...)
}

func (r *repository) GetAllImageByDataset(dID uint64) ([]Image, error) {
	imgs := make([]Image, 0)
	k := rediskey.ImageAllByDatasetID(dID)
	err := r.cacheRepo.Get(k, &imgs)
	if err == nil {
		logger.Infof("cache hit getting all images of project %d", dID)
		return imgs, nil
	}
	if errors.Type(err) == errors.CacheNotFound {
		logger.Infof("cache miss getting all images of project %d", dID)
	} else {
		logger.Error("cannot get all images of dataset %d from cache. error %v", dID, err)
	}
	imgs, err = r.dbRepo.GetAllByDataset(dID)
	if err != nil {
		return nil, err
	}
	logger.Infof("get all images from cache successfully")
	return imgs, nil
}

func (r *repository) BulkInsert(images []Image, dID uint64) error {
	return r.dbRepo.BulkInsert(images, dID)
}
