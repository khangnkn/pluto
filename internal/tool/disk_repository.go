package tool

import (
	"github.com/jinzhu/gorm"
	"github.com/nkhang/pluto/pkg/errors"
)

type diskRepository struct {
	db *gorm.DB
}

func NewDiskRepository(db *gorm.DB) *diskRepository {
	return &diskRepository{
		db: db,
	}
}

func (d *diskRepository) GetAll() ([]Tool, error) {
	t := make([]Tool, 0)
	err := d.db.Model(&Tool{}).Find(&t).Error
	if len(t) == 0 {
		return nil, errors.ToolNoRecord.NewWithMessage("no record found")
	}
	if err != nil {
		return nil, errors.ToolQueryError.Wrap(err, "cannot query all tools")
	}
	return t, nil
}
