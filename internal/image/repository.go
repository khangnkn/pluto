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
	CreateImage(title, url, thumbnail string, w, h int, size int64, dataset_id uint64) (Image, error)
	Incr(id uint64) error
	BulkInsert(images []Image, dID uint64) error
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

func (r *repository) CreateImage(title, url, thumbnail string, w, h int, size int64, datasetId uint64) (Image, error) {
	r.InvalidateDatasetImage(datasetId)
	return r.dbRepo.CreateImage(title, url, thumbnail, w, h, size, datasetId)
}

func (r *repository) InvalidateDatasetImage(dID uint64) {
	pattern := rediskey.ImageByDatasetIDAllKeys(dID)
	keys, err := r.cacheRepo.Keys(pattern)
	if len(keys) == 0 {
		logger.Info("[IMAGE] - no keys found to invalidate")
		return
	}
	if err != nil {
		logger.Errorf("[IMAGE] - error getting all keys from redis. err %v", err)
		return
	}
	logger.Infof("the following keys will be deleted: %v", keys)
	err = r.cacheRepo.Del(keys...)
	if err != nil {
		logger.Errorf("[IMAGE] - error deleting keys. err %v", err)
	}
}

func (r *repository) GetAllImageByDataset(dID uint64) ([]Image, error) {
	images := make([]Image, 0)
	k := rediskey.ImageAllByDatasetID(dID)
	err := r.cacheRepo.Get(k, &images)
	if err == nil {
		logger.Infof("[IMAGE] - cache hit getting all images of project %d", dID)
		return images, nil
	}
	if errors.Type(err) == errors.CacheNotFound {
		logger.Infof("[IMAGE] - cache miss getting all images of project %d", dID)
	} else {
		logger.Error("[IMAGE] - cannot get all images of dataset %d from cache. error %v", dID, err)
	}
	images, err = r.dbRepo.GetAllByDataset(dID)
	if err != nil {
		return nil, err
	}
	logger.Infof("get all images from cache successfully")
	return images, nil
}

func (r *repository) BulkInsert(images []Image, dID uint64) error {
	r.InvalidateDatasetImage(dID)
	return r.dbRepo.BulkInsert(images, dID)
}

func (r *repository) Incr(id uint64) error {
	r.InvalidateDatasetImage(id)
	return r.dbRepo.Incr(id)
}
