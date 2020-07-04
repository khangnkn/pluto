package task

import (
	"github.com/jinzhu/gorm"
	"github.com/nkhang/pluto/pkg/errors"
	gormbulk "github.com/t-tiger/gorm-bulk-insert/v2"
)

type DBRepository interface {
	CreateTask(assigner, labeler, reviewer, datasetID uint64) (Task, error)
	AddImages(id uint64, imageIDs []uint64) error
	GetTaskDetails(taskID uint64, offset, limit int) ([]Detail, error)
}

type dbRepository struct {
	db *gorm.DB
}

func NewDBRepository(db *gorm.DB) *dbRepository {
	return &dbRepository{db: db}
}

func (r *dbRepository) CreateTask(assigner, labeler, reviewer, datasetID uint64) (Task, error) {
	t := Task{
		DatasetID: datasetID,
		Assigner:  assigner,
		Labeler:   labeler,
		Reviewer:  reviewer,
	}
	err := r.db.Create(&t).Error
	if err != nil {
		return Task{}, errors.TaskCannotCreate.Wrap(err, "cannot create task")
	}
	return t, nil
}

func (r *dbRepository) AddImages(id uint64, imageIDs []uint64) error {
	records := make([]interface{}, len(imageIDs))
	for i := range records {
		var record = Detail{
			Status:  Unassigned,
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

func (r *dbRepository) GetTaskDetails(taskID uint64, offset, limit int) ([]Detail, error) {
	var details []Detail
	var tableName = Detail{TaskID: taskID}.TableName()
	err := r.db.Table(tableName).
		Preload("Image").
		Where("task_id = ?", taskID).
		Offset(offset).
		Limit(limit).
		Find(&details).Error
	if err != nil {
		return nil, errors.TaskDetailCannotGet.NewWithMessage("cannot get task details")
	}
	return details, nil
}
