package datasetfx

import (
	"github.com/jinzhu/gorm"

	"github.com/nkhang/pluto/internal/dataset"
	"github.com/nkhang/pluto/internal/dataset/datasetapi"
	"github.com/nkhang/pluto/internal/image"
	"github.com/nkhang/pluto/pkg/cache"
	"github.com/nkhang/pluto/pkg/pgin"
)

func provideRepository(db *gorm.DB, c cache.Cache) dataset.Repository {
	dbRepo := dataset.NewDbRepository(db)
	return dataset.NewRepository(dbRepo, c)
}

func provideAPIRepo(r dataset.Repository, imgRepo image.Repository) datasetapi.Repository {
	return datasetapi.NewRepository(r, imgRepo)
}

func provideService(repository datasetapi.Repository) pgin.StandaloneRouter {
	return datasetapi.NewService(repository)
}
