package tool

import (
	"github.com/jinzhu/gorm"
	"github.com/nkhang/pluto/pkg/errors"
)

type DbRepository interface {
	GetAll() ([]Tool, error)
}

type dbRepository struct {
	db *gorm.DB
}

func NewDiskRepository(db *gorm.DB) *dbRepository {
	return &dbRepository{
		db: db,
	}
}

func (d *dbRepository) GetAll() ([]Tool, error) {
	t := make([]Tool, 0)
	err := d.db.Model(&Tool{}).Find(&t).Error
	if err != nil {
		return nil, errors.ToolQueryError.Wrap(err, "cannot query all tools")
	}
	return t, nil
}
