package tool

import "github.com/nkhang/pluto/pkg/gorm"

type Tool struct {
	gorm.Model
	Name string
}
