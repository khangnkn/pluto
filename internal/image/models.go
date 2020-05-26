package image

import "github.com/nkhang/pluto/pkg/gorm"

type Image struct {
	gorm.Model
	URL       string
	Title     string
	Width     int
	Height    int
	Size      int64
	DatasetID uint64
}
