package image

import (
	"github.com/jinzhu/gorm"
	gormbulk "github.com/t-tiger/gorm-bulk-insert/v2"

	"github.com/nkhang/pluto/pkg/errors"
)

type DBRepository interface {
	Get(id uint64) (Image, error)
	GetByDataset(dID uint64, offset, limit int) (imgs []Image, err error)
	CreateImage(title, url string, w, h int, size int64, dataset_id uint64) (Image, error)
	GetAllByDataset(dID uint64) (images []Image, err error)
	BulkInsert(images []Image, dID uint64) error
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

func (r *dbRepository) GetAllByDataset(dID uint64) (images []Image, err error) {
	err = r.db.Where("dataset_id = ?", dID).
		Find(&images).
		Order("status asc").Error
	if err != nil {
		err = errors.ImageQueryError.Wrap(err, "images query error")
		return
	}
	return
}

func (r *dbRepository) CreateImage(title, url string, w, h int, size int64, dID uint64) (Image, error) {
	img := Image{
		URL:       url,
		Title:     title,
		Width:     w,
		Height:    h,
		Size:      size,
		DatasetID: dID,
	}
	err := r.db.Save(&img).Error
	if err != nil {
		return Image{}, errors.ImageErrorCreating.NewWithMessage("error creating image")
	}
	return img, nil
}

func (r *dbRepository) BulkInsert(images []Image, dID uint64) error {
	var clone []interface{}
	for i := range images {
		var img = Image{
			Title:     images[i].Title,
			URL:       images[i].URL,
			Width:     images[i].Width,
			Height:    images[i].Height,
			Size:      images[i].Size,
			DatasetID: dID,
		}
		clone = append(clone, img)
	}
	err := gormbulk.BulkInsert(r.db, clone, 1000)
	if err != nil {
		return errors.ImageErrorBulkCreating.Wrap(err, "cannot creating bulk image")
	}
	return nil
}
