package label

import (
	"github.com/nkhang/pluto/internal/rediskey"
	"github.com/nkhang/pluto/pkg/cache"
	"github.com/nkhang/pluto/pkg/errors"
	"github.com/nkhang/pluto/pkg/logger"
)

type Repository interface {
	GetByProjectId(pID uint64) ([]Label, error)
	CreateLabel(name, color string, projectID, toolID uint64) error
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

func (r *repository) GetByProjectId(pID uint64) ([]Label, error) {
	var labels = make([]Label, 0)
	k := rediskey.LabelsByProject(pID)
	err := r.cacheRepo.Get(k, &labels)
	if err == nil {
		logger.Infof("cache hit for labels by project [%d]", pID)
		return labels, nil
	}
	if errors.Type(err) == errors.CacheNotFound {
		logger.Infof("cache miss for labels by project [%d]", pID)
	} else {
		logger.Infof("errors getting labels for project [%d]", pID)
	}
	labels, err = r.dbRepo.GetByProjectID(pID)
	if err != nil {
		return nil, err
	}
	go func() {
		err := r.cacheRepo.Set(k, &labels)
		if err != nil {
			logger.Error(err)
		}
	}()
	return labels, nil
}
func (r *repository) CreateLabel(name, color string, projectID, toolID uint64) error {
	k := rediskey.LabelsByProject(projectID)
	go func() {
		if err := r.cacheRepo.Del(k); err != nil {
			logger.Errorf("cannot invalidate all tools for project %d, error %s", projectID, err.Error())
		}
	}()
	return r.dbRepo.CreateLabel(name, color, projectID, toolID)
}
