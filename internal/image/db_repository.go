package image

import (
	"github.com/jinzhu/gorm"

	"github.com/nkhang/pluto/pkg/errors"
)

type DBRepository interface {
	Get(id uint64) (Image, error)
	GetByDataset(dID uint64, offset, limit int) (imgs []Image, err error)
	CreateImage(name string, w, h int, dataset_id uint64) (Image, error)
}

type dbRepository struct {
	db *gorm.DB
}

func NewDBRepository(db *gorm.DB) *dbRepository {
	return &dbRepository{db: db}
}

func (r *dbRepository) Get(id uint64) (img Image, err error) {
	result := r.db.First(&img, id)
	if result.RecordNotFound() {
		err = errors.ImageNotFound.NewWithMessage("image not found")
		return
	}
	if err = result.Error; err != nil {
		err = errors.ImageQueryError.Wrap(err, "image query error")
		return
	}
	return
}

func (r *dbRepository) GetByDataset(dID uint64, offset, limit int) (images []Image, err error) {
	err = r.db.Where("dataset_id = ?", dID).
		Offset(offset).
		Limit(limit).
		Find(&images).Error
	if err != nil {
		err = errors.ImageQueryError.Wrap(err, "images query error")
		return
	}
	return
}

func (r *dbRepository) CreateImage(name string, w, h int, dataset_id uint64) (Image, error) {
	img := Image{
		URL:       name,
		Width:     w,
		Height:    h,
		DatasetID: dataset_id,
	}
	err := r.db.Save(&img).Error
	if err != nil {
		return Image{}, errors.ImageErrorCreating.NewWithMessage("error creating image")
	}
	return img, nil
}
