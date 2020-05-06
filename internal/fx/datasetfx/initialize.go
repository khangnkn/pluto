package datasetfx

import (
	"github.com/jinzhu/gorm"

	"github.com/nkhang/pluto/internal/dataset"
	"github.com/nkhang/pluto/internal/dataset/datasetapi"
	"github.com/nkhang/pluto/pkg/cache"
	"github.com/nkhang/pluto/pkg/gin"
)

func provideRepository(db *gorm.DB, c cache.Cache) dataset.Repository {
	dbRepo := dataset.NewDbRepository(db)
	return dataset.NewRepository(dbRepo, c)
}

func provideService(r dataset.Repository) gin.IEngine {
	repository := datasetapi.NewRepository(r)
	return datasetapi.NewService(repository)
}
