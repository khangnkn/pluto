package task

import (
	"github.com/nkhang/pluto/internal/rediskey"
	"github.com/nkhang/pluto/pkg/cache"
	"github.com/nkhang/pluto/pkg/errors"
	"github.com/nkhang/pluto/pkg/logger"
)

type Repository interface {
	GetTask(taskID uint64) (Task, error)
	CreateTask(title, description string, assigner, labeler, reviewer, projectID, datasetID uint64, images []uint64) (Task, error)
	GetTasksByUser(userID uint64, role Role, status Status, offset, limit int) (tasks []Task, total int, err error)
	GetTasksByProject(projectID uint64, status Status, offset, limit int) (tasks []Task, total int, err error)
	GetByProjectAndUser(projectID, userID uint64, role Role, offset, limit int) (tasks []Task, total int, err error)
	DeleteTask(taskID uint64) error
	DeleteTaskByProject(projectID uint64) error
	GetTaskDetails(taskID uint64, status DetailStatus, currentID uint64, limit int) ([]Detail, int, error)
	UpdateTask(taskID uint64, changes map[string]interface{}) (Task, error)
	UpdateTaskDetail(taskID, detailID uint64, changes map[string]interface{}) (Detail, error)
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

func (r *repository) GetTask(taskID uint64) (task Task, err error) {
	key := rediskey.TaskByID(taskID)
	err = r.cache.Get(key, &task)
	if err == nil {
		return
	}
	task, err = r.dbRepo.GetTask(taskID)
	if err != nil {
		return
	}
	go func() {
		r.cache.Set(key, &task)
	}()
	return
}

func (r *repository) CreateTask(title, description string, assigner, labeler, reviewer, projectID, datasetID uint64, images []uint64) (Task, error) {
	task, err := r.dbRepo.CreateTask(title, description, assigner, labeler, reviewer, projectID, datasetID)
	if err != nil {
		return Task{}, err
	}
	go func() {
		r.InvalidateForUser(labeler)
		r.InvalidateForUser(reviewer)
		_, _, pattern := rediskey.TaskByProject(projectID, 0, 0, 0)
		keys, err := r.cache.Keys(pattern)
		if err != nil {
			logger.Errorf("error getting keys for pattern %d. error %v", pattern, err)
			return
		}
		if err := r.cache.Del(keys...); err != nil {
			logger.Error(err)
		}
	}()
	if err := r.dbRepo.AddImages(task.ID, images); err != nil {
		return Task{}, errors.TaskCannotCreate.NewWithMessage("task created, but images not add properly")
	}
	return task, nil
}

func (r *repository) GetTasksByUser(userID uint64, role Role, status Status, offset, limit int) (tasks []Task, total int, err error) {
	tasks, total, err = r.dbRepo.GetTasksByUser(userID, role, status, offset, limit)
	return
}

func (r *repository) GetTasksByProject(projectID uint64, status Status, offset, limit int) (tasks []Task, total int, err error) {
	specificKey, totalKey, _ := rediskey.TaskByProject(projectID, uint32(status), offset, limit)
	err1 := r.cache.Get(specificKey, &tasks)
	err2 := r.cache.Get(totalKey, &total)
	if err1 == nil && err2 == nil {
		logger.Infof("cache hit getting all tasks by projects %d", projectID)
		return
	}
	tasks, total, err = r.dbRepo.GetTasksByProject(projectID, status, offset, limit)
	if err != nil {
		return
	}
	go func() {
		r.cache.Set(specificKey, &tasks)
		r.cache.Set(totalKey, &total)
	}()
	return
}

func (r *repository) GetByProjectAndUser(projectID, userID uint64, role Role, offset, limit int) (tasks []Task, total int, err error) {
	return r.dbRepo.GetByProjectAndUser(projectID, userID, role, offset, limit)
}

func (r *repository) DeleteTask(id uint64) error {
	task, err := r.GetTask(id)
	if err != nil {
		return err
	}
	go func() {
		r.InvalidateForUser(task.Reviewer)
		r.InvalidateForUser(task.Labeler)
		r.InvalidateForProject(task.ProjectID)
		k := rediskey.TaskByID(id)
		r.cache.Del(k)
	}()
	err = r.dbRepo.DeleteTask(id)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) GetTaskDetails(taskID uint64, status DetailStatus, currentID uint64, limit int) ([]Detail, int, error) {
	details, total, count, err := r.dbRepo.GetTaskDetails(taskID, status, currentID, limit)
	if err != nil {
		return nil, 0, err
	}
	logger.Infof("get images of tasks %d, status %d return %d image", taskID, status, count)
	if count == 0 {
		var val int32
		switch status {
		case 0:
			val = 2
		case 2:
			val = 3
		default:
			val = 3
		}
		_, err = r.UpdateTask(taskID, map[string]interface{}{"status": val})
		if err != nil {
			return nil, 0, err
		}
	}
	return details, total, nil
}

func (r *repository) UpdateTaskDetail(taskID, detailID uint64, changes map[string]interface{}) (Detail, error) {
	return r.dbRepo.UpdateTaskDetail(taskID, detailID, changes)
}

func (r *repository) InvalidateForUser(userID uint64) {
	_, _, pattern := rediskey.TaskByUser(userID, 0, 0, 0, 0)
	keys, err := r.cache.Keys(pattern)
	if err != nil {
		logger.Errorf("error get pattern %s", pattern)
		return
	}
	err = r.cache.Del(keys...)
	if err != nil {
		logger.Errorf("error delete keys %v", keys)
	}
}
func (r *repository) InvalidateForProject(projectID uint64) {
	_, _, pattern := rediskey.TaskByProject(projectID, 0, 0, 0)
	keys, err := r.cache.Keys(pattern)
	if err != nil {
		logger.Errorf("error get pattern %s", pattern)
		return
	}
	err = r.cache.Del(keys...)
	if err != nil {
		logger.Errorf("error delete keys %v", keys)
	}
}

func (r *repository) DeleteTaskByProject(projectID uint64) error {
	err := r.dbRepo.DeleteTaskByProject(projectID)
	if err != nil {
		return err
	}
	r.InvalidateForProject(projectID)
	return nil
}

func (r *repository) UpdateTask(taskID uint64, changes map[string]interface{}) (Task, error) {
	return r.dbRepo.UpdateTask(taskID, changes)
}

func GetUserTaskInProject() {

}
