package image

import "github.com/nkhang/pluto/pkg/gorm"

type Image struct {
	gorm.Model
	URL       string
	Width     uint32
	Height    uint32
	DatasetID uint64
}
