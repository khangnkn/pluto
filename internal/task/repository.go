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
	CheckTaskStatus(taskID uint64, detailStatus DetailStatus) error
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
		r.invalidateForUser(labeler)
		r.invalidateForUser(reviewer)
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
		r.invalidateForUser(task.Reviewer)
		r.invalidateForUser(task.Labeler)
		r.invalidateForUser(task.Assigner)
		r.invalidateForProject(task.ProjectID)
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
	details, total, err := r.dbRepo.GetTaskDetails(taskID, status, currentID, limit)
	if err != nil {
		return nil, 0, err
	}
	return details, total, nil
}

func (r *repository) UpdateTaskDetail(taskID, detailID uint64, changes map[string]interface{}) (Detail, error) {
	return r.dbRepo.UpdateTaskDetail(taskID, detailID, changes)
}

func (r *repository) invalidateForUser(userID uint64) {
	_, _, pattern := rediskey.TaskByUser(userID, 0, 0, 0, 0)
	keys, err := r.cache.Keys(pattern)
	if err != nil {
		logger.Errorf("error get pattern %s", pattern)
		return
	}
	err = r.cache.Del(keys...)
	if err != nil {
		logger.Errorf("error delete keys %v", keys)
		return
	}
	logger.Infof("[TASK] - invalidate tasks for user %d successfully", userID)
}
func (r *repository) invalidateForProject(projectID uint64) {
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
	logger.Infof("[TASK] - invalidate tasks for project %d successfully", projectID)
}

func (r *repository) DeleteTaskByProject(projectID uint64) error {
	err := r.dbRepo.DeleteTaskByProject(projectID)
	if err != nil {
		return err
	}
	r.invalidateForProject(projectID)
	return nil
}

func (r *repository) UpdateTask(taskID uint64, changes map[string]interface{}) (Task, error) {
	task, err := r.dbRepo.UpdateTask(taskID, changes)
	if err != nil {
		return Task{}, err
	}
	r.invalidateForProject(task.ProjectID)
	return task, nil
}

func (r *repository) CheckTaskStatus(taskID uint64, detailStatus DetailStatus) (err error) {
	rl, s := relative(detailStatus)
	var (
		buffer []Detail
		images = make([]Detail, 0)
		total  int
	)
	for i := range rl {
		buffer, total, err = r.dbRepo.GetTaskDetails(taskID, rl[i], 0, 0)
		if err != nil {
			return err
		}
		images = append(images, buffer...)
	}
	if len(images) != total {
		return nil
	}
	logger.Infof("[TASK] no more images of status %v of task %d - increasing task status...", rl, taskID)
	_, err = r.UpdateTask(taskID, map[string]interface{}{
		"status": s,
	})
	return
}
