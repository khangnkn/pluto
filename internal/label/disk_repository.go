package label

import (
	"fmt"

	"github.com/jinzhu/gorm"

	"github.com/nkhang/pluto/pkg/errors"
)

type DBRepository interface {
	GetByProjectID(projectID uint64) ([]Label, error)
}

type dbRepository struct {
	db *gorm.DB
}

func NewDiskRepository(db *gorm.DB) *dbRepository {
	return &dbRepository{db: db}
}

func (d *dbRepository) GetByProjectID(projectID uint64) ([]Label, error) {
	l := make([]Label, 0)
	query := fmt.Sprint(fieldProjectID, " = ?")
	err := d.db.Where(query, projectID).Preload("Tool").Find(&l).Error
	if err != nil {
		return nil, errors.LabelQueryError.NewWithMessage("label query error")
	}
	return l, nil
}
