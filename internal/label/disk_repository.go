package label

import (
	"fmt"

	"github.com/jinzhu/gorm"

	"github.com/nkhang/pluto/pkg/errors"
)

type dbRepository interface {
	GetByProjectID(projectID uint64) ([]Label, error)
}

type diskRepository struct {
	db *gorm.DB
}

func NewDiskRepository(db *gorm.DB) *diskRepository {
	return &diskRepository{db: db}
}

func (d *diskRepository) GetByProjectID(projectID uint64) ([]Label, error) {
	l := make([]Label, 0)
	query := fmt.Sprint(fieldProjectID, " = ?")
	err := d.db.Where(query, projectID).Preload("Tool").Find(&l).Error
	if err != nil {
		return nil, errors.LabelQueryError.NewWithMessage("label query error")
	}
	if len(l) == 0 {
		return nil, errors.LabelRecordNotFound.NewWithMessage("label not found")
	}
	return l, nil
}
