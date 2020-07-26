package dataset

import (
	"github.com/nkhang/pluto/internal/rediskey"
	"github.com/nkhang/pluto/internal/task"
	"github.com/nkhang/pluto/pkg/cache"
	"github.com/nkhang/pluto/pkg/errors"
	"github.com/nkhang/pluto/pkg/logger"
)

type Repository DbRepository

type repository struct {
	dbRepo    DbRepository
	taskRepo  task.Repository
	cacheRepo cache.Cache
}

func NewRepository(d DbRepository, c cache.Cache, t task.Repository) *repository {
	return &repository{
		dbRepo:    d,
		cacheRepo: c,
		taskRepo:  t,
	}
}

func (r *repository) Get(dID uint64) (d Dataset, err error) {
	k := rediskey.DatasetByID(dID)
	err = r.cacheRepo.Get(k, &d)
	if err == nil {
		logger.Infof("cache hit getting dataset %d", dID)
		return d, nil
	}
	if errors.Type(err) == errors.CacheNotFound {
		logger.Infof("cache miss getting dataset %d", dID)
	} else {
		logger.Errorf("error getting dataset %d from cache", dID)
	}
	d, err = r.dbRepo.Get(dID)
	if err != nil {
		return Dataset{}, err
	}
	go func() {
		err := r.cacheRepo.Set(k, &d)
		if err != nil {
			logger.Error("error set dataset %d to cache", dID)
		}
	}()
	return d, nil
}

func (r *repository) GetByProject(pID uint64) ([]Dataset, error) {
	var ds = make([]Dataset, 0)
	k := rediskey.DatasetByProject(pID)
	err := r.cacheRepo.Get(k, &ds)
	if err == nil {
		logger.Infof("cache hit getting datasets of project %d", pID)
		return ds, nil
	}
	if errors.Type(err) == errors.CacheNotFound {
		logger.Infof("cache miss getting datasets of project %d", pID)
	} else {
		logger.Errorf("error getting datasets of project %d from cache", pID)
	}
	ds, err = r.dbRepo.GetByProject(pID)
	if err != nil {
		logger.Error("error getting datasets of projects %d from database", pID)
		return nil, err
	}
	go func() {
		err := r.cacheRepo.Set(k, &ds)
		if err != nil {
			logger.Error("error set datasets of projects %d to cache", pID)
		}
	}()
	return ds, nil
}

func (r *repository) CreateDataset(title, description string, pID uint64) (Dataset, error) {
	go func() {
		k := rediskey.DatasetByProject(pID)
		err := r.cacheRepo.Del(k)
		if err != nil {
			logger.Infof("cannot invalidate cache for all dataset by project %d", pID)
		}
		logger.Infof("invalidate cache for project %d", pID)
	}()
	return r.dbRepo.CreateDataset(title, description, pID)
}

func (r *repository) DeleteDataset(ID uint64) error {
	d, err := r.Get(ID)
	if err != nil {
		return err
	}
	go func() {
		r.invalidate(ID, d.ProjectID)
	}()
	tasks, _, err := r.taskRepo.GetTasksByProject(d.ProjectID, task.Any, 0, 0)
	if err != nil {
		return err
	}
	for _, t := range tasks {
		if t.DatasetID == ID {
			err := r.taskRepo.DeleteTask(t.ID)
			if err != nil {
				logger.Errorf("[DATASET] - error delete task %d of dataset %d", t.ID, ID)
			}
		}
	}
	err = r.dbRepo.DeleteDataset(ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) invalidate(datasetID, projectID uint64) {
	k := rediskey.DatasetByID(datasetID)
	k2 := rediskey.DatasetByProject(projectID)
	err := r.cacheRepo.Del(k, k2)
	if err != nil {
		logger.Errorf("cannot delete dataset %d from cache", datasetID)
	}
}

func (r *repository) Update(id uint64, changes map[string]interface{}) (Dataset, error) {
	d, err := r.dbRepo.Update(id, changes)
	if err != nil {
		return Dataset{}, err
	}
	r.invalidate(d.ID, d.ProjectID)
	return d, nil
}

func (r *repository) DeleteByProject(projectID uint64) error {
	err := r.dbRepo.DeleteByProject(projectID)
	if err != nil {
		return err
	}
	r.invalidate(0, projectID)
	return nil
}
