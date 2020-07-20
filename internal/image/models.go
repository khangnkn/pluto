package image

import "github.com/nkhang/pluto/pkg/gorm"

type Image struct {
	gorm.Model
	URL       string
	Status    uint32
	Title     string
	Width     int
	Height    int
	Size      int64
	DatasetID uint64
}
