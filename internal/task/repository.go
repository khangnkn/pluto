package task

import "github.com/nkhang/pluto/pkg/cache"

type Repository interface {
	CreateTask(assigner, labeler, reviewer, datasetID uint64, images []uint64) error
	//GetImages(id uint64)
}

type repository struct {
	dbRepo DBRepository
	cache  cache.Cache
}

func NewRepository(dbRepo DBRepository, cache cache.Cache) *repository {
	return &repository{
		dbRepo: dbRepo,
		cache:  cache,
	}
}

func (r *repository) CreateTask(assigner, labeler, reviewer, datasetID uint64, images []uint64) error {
	task, err := r.dbRepo.CreateTask(assigner, labeler, reviewer, datasetID)
	if err != nil {
		return err
	}
	return r.dbRepo.AddImages(task.ID, images)
}
