package label

import (
	"fmt"

	"github.com/nkhang/pluto/pkg/logger"

	"github.com/jinzhu/gorm"

	"github.com/nkhang/pluto/pkg/errors"
)

type DBRepository interface {
	GetByProjectID(projectID uint64) ([]Label, error)
	CreateLabel(name, color string, projectID, toolID uint64) error
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
	err := d.db.Preload("Tool").Where(query, projectID).Find(&l).Error
	logger.Infof("%+v", l)
	if err != nil {
		return nil, errors.LabelQueryError.NewWithMessage("label query error")
	}
	return l, nil
}

func (d *dbRepository) CreateLabel(name, color string, projectID, toolID uint64) error {
	l := Label{
		Name:      name,
		Color:     color,
		ProjectID: projectID,
		ToolID:    toolID,
	}
	err := d.db.Create(&l).Error
	if err != nil {
		return errors.LabelCannotCreate.Wrap(err, "cannot create label")
	}
	return nil
}
