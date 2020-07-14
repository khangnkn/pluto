package task

import (
	"github.com/jinzhu/gorm"
	"github.com/nkhang/pluto/internal/dataset"
	"github.com/nkhang/pluto/internal/project"
	"github.com/nkhang/pluto/pkg/errors"
	gormbulk "github.com/t-tiger/gorm-bulk-insert/v2"
)

type DBRepository interface {
	GetTask(taskID uint64) (Task, error)
	GetTasksByProject(projectID uint64, status Status, offset, limit int) (tasks []Task, total int, err error)
	GetTasksByUser(userID uint64, role Role, status Status, offset, limit int) (tasks []Task, total int, err error)
	CreateTask(title, description string, assigner, labeler, reviewer, projectID, datasetID uint64) (Task, error)
	DeleteTask(id uint64) error
	DeleteTaskByProject(projectID uint64) error
	AddImages(id uint64, imageIDs []uint64) error
	GetTaskDetails(taskID uint64, currentID uint64, limit int) ([]Detail, error)
	UpdateTaskDetail(taskID, detailID uint64, changes map[string]interface{}) (Detail, error)
}

type dbRepository struct {
	db *gorm.DB
}

func NewDBRepository(db *gorm.DB) *dbRepository {
	return &dbRepository{db: db}
}

func (r *dbRepository) GetTask(taskID uint64) (task Task, err error) {
	db := r.db.First(&task, taskID)
	if db.RecordNotFound() {
		return Task{}, errors.TaskNotFound.NewWithMessageF("task %d not found", taskID)
	}
	if err = db.Error; err != nil {
		err = errors.TaskCannotGet.Wrap(err, "cannot get task")
		return
	}
	return
}

func (r *dbRepository) GetTasksByUser(userID uint64, role Role, status Status, offset, limit int) (tasks []Task, total int, err error) {
	db := r.db.Model(&Task{})
	switch role {
	case AnyRole:
		db = db.Where(&Task{Labeler: userID}).Or(&Task{Reviewer: userID})
	case Labeler:
		db = db.Where(&Task{Labeler: userID})
	case Reviewer:
		db = db.Where(&Task{Reviewer: userID})
	default:
		err = errors.TaskCannotGet.NewWithMessage("role not supported")
		return
	}
	if status != Any {
		db.Where("status = ?", status)
	}
	db = db.Count(&total)
	if offset != 0 || limit != 0 {
		db = db.Offset(offset).Limit(limit)
	}
	err = db.Find(&tasks).Error
	if err != nil {
		return nil, 0, errors.TaskCannotGet.Wrap(err, "cannot get task")
	}
	return
}

func (r *dbRepository) GetTasksByProject(projectID uint64, status Status, offset, limit int) (tasks []Task, total int, err error) {
	db := r.db.Model(&Task{}).Where(&Task{ProjectID: projectID})
	if status != Any {
		db.Where("status = ?", status)
	}
	db = db.Count(&total)
	if offset != 0 || limit != 0 {
		db = db.Offset(offset).Limit(limit)
	}
	err = db.Find(&tasks).Error
	if err != nil {
		return nil, 0, errors.TaskCannotGet.Wrap(err, "cannot get task")
	}
	return
}

func (r *dbRepository) CreateTask(title, description string, assigner, labeler, reviewer, projectID, datasetID uint64) (Task, error) {
	t := Task{
		Title:       title,
		Description: description,
		ProjectID:   projectID,
		DatasetID:   datasetID,
		Assigner:    assigner,
		Labeler:     labeler,
		Reviewer:    reviewer,
		Status:      Labeling,
	}
	if r.db.First(&project.Project{}, projectID).RecordNotFound() {
		return Task{}, errors.TaskCannotCreate.NewWithMessage("project not found")
	}
	var d dataset.Dataset
	if r.db.First(&d, datasetID).RecordNotFound() || d.ProjectID != projectID {
		return Task{}, errors.TaskCannotCreate.NewWithMessage("dataset not found")
	}
	err := r.db.Create(&t).Error
	if err != nil {
		return Task{}, errors.TaskCannotCreate.Wrap(err, "cannot create task")
	}
	return t, nil
}

func (r *dbRepository) DeleteTask(id uint64) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var t Task
		var d Detail
		d.TaskID = id
		t.ID = id
		err := tx.Model(&t).Delete(&Task{}, id).Error
		if err != nil {
			return err
		}
		err = tx.Model(&d).
			Table(d.TableName()).
			Where("task_id = ?", id).
			Delete(&Detail{}).Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return errors.TaskDetailCannotDelete.Wrap(err, "cannot delete task")
	}
	return nil
}

func (r *dbRepository) AddImages(id uint64, imageIDs []uint64) error {
	records := make([]interface{}, len(imageIDs))
	for i := range records {
		var record = Detail{
			Status:  Pending,
			TaskID:  id,
			ImageID: imageIDs[i],
		}
		records[i] = record
	}
	err := gormbulk.BulkInsert(r.db, records, 1000)
	if err != nil {
		return errors.TaskCannotCreate.Wrap(err, "cannot create tasks")
	}
	return nil
}

func (r *dbRepository) GetTaskDetails(taskID uint64, currentID uint64, limit int) ([]Detail, error) {
	var details []Detail
	var tableName = Detail{TaskID: taskID}.TableName()
	db := r.db.Table(tableName).
		Preload("Image").
		Where("task_id = ? and status = ?", taskID, Pending)
	if currentID != 0 || limit != 0 {
		db = db.Where("id > ?", currentID).Limit(limit)
	}
	err := db.Find(&details).Error
	if err != nil {
		return nil, errors.TaskDetailCannotGet.NewWithMessage("cannot get task details")
	}
	return details, nil
}

func (r *dbRepository) UpdateTaskDetail(taskID, detailID uint64, changes map[string]interface{}) (Detail, error) {
	var detail = Detail{TaskID: taskID}
	detail.ID = detailID
	err := r.db.Model(&detail).Update(changes).Preload("Image").First(&detail, detailID).Error
	if err != nil {
		return Detail{}, errors.TaskDetailCannotUpdate.Wrap(err, "cannot update task detail")
	}
	return detail, nil
}

func (r *dbRepository) DeleteTaskByProject(projectID uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var tasks = make([]Task, 0)
		err := r.db.Where("project_id = ?", projectID).Find(&tasks).Delete(&Task{}).Error
		if err != nil {
			return errors.TaskCannotDelete.NewWithMessageF("cannot delete tasks of project %d", projectID)
		}
		for _, task := range tasks {
			err = r.db.Where("task_id = ?", task.ID).Delete(&Detail{}).Error
			if err != nil {
				return errors.TaskDetailCannotDelete.Wrap(err, "cannot delete task detail")
			}
		}
		return nil
	})
}
