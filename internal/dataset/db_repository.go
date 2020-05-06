package dataset

import (
	"github.com/jinzhu/gorm"

	"github.com/nkhang/pluto/pkg/errors"
)

type DbRepository interface {
	Get(dID uint64) (Dataset, error)
	GetByProject(pID uint64) ([]Dataset, error)
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
	result := r.db.First(&d, dID)
	if result.RecordNotFound() {
		err = errors.DatasetNotFound.NewWithMessage("dataset not found")
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
