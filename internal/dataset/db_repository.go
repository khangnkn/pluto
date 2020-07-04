package dataset

import (
	"fmt"

	"github.com/jinzhu/gorm"

	"github.com/nkhang/pluto/pkg/errors"
)

type DbRepository interface {
	Get(dID uint64) (Dataset, error)
	GetByProject(pID uint64) ([]Dataset, error)
	CreateDataset(title, description string, pID uint64) (Dataset, error)
	DeleteDataset(ID uint64) error
}

type dbRepository struct {
	db *gorm.DB
}

func NewDbRepository(db *gorm.DB) *dbRepository {
	return &dbRepository{
		db: db,
	}
}

func (r *dbRepository) Get(dID uint64) (d Dataset, err error) {
	result := r.db.Preload("Project").First(&d, dID)
	if result.RecordNotFound() {
		err = errors.DatasetNotFound.NewWithMessage("dataset not found")
		return
	}
	if err = result.Error; err != nil {
		err = errors.DatasetQueryError.Wrap(err, "dataset query error")
		return
	}
	return d, nil
}

func (r *dbRepository) GetByProject(pID uint64) ([]Dataset, error) {
	result := make([]Dataset, 0)
	err := r.db.Where(fieldProjectID+" = ?", pID).Find(&result).Error
	if err != nil {
		return nil, errors.DatasetQueryError.Wrap(err, "dataset query error")
	}
	return result, nil
}

func (r *dbRepository) CreateDataset(title, description string, pID uint64) (Dataset, error) {
	d := Dataset{
		Title:       title,
		Description: description,
		ProjectID:   pID,
	}
	err := r.db.Create(&d).Error
	if err != nil {
		return Dataset{}, errors.DatasetCannotCreate.Wrap(err, "cannot create dataset")
	}
	return d, nil
}

func (r *dbRepository) DeleteDataset(ID uint64) error {
	err := r.db.Delete(&Dataset{}, ID).Error
	if err != nil {
		return errors.DatasetCannotDelete.Wrap(err, fmt.Sprintf("cannot delete dataset %d", ID))
	}
	return nil
}
