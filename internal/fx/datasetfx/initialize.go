package datasetfx

import (
	"github.com/jinzhu/gorm"
	"github.com/nkhang/pluto/internal/project"
	"github.com/nkhang/pluto/internal/task"
	"go.uber.org/fx"

	"github.com/nkhang/pluto/internal/dataset"
	"github.com/nkhang/pluto/internal/dataset/datasetapi"
	"github.com/nkhang/pluto/internal/image"
	"github.com/nkhang/pluto/pkg/cache"
	"github.com/nkhang/pluto/pkg/pgin"
)

func provideRepository(db *gorm.DB, c cache.Cache, t task.Repository) dataset.Repository {
	dbRepo := dataset.NewDbRepository(db)
	return dataset.NewRepository(dbRepo, c, t)
}

func provideAPIRepo(r dataset.Repository, imgRepo image.Repository, p project.Repository) datasetapi.Repository {
	return datasetapi.NewRepository(r, imgRepo, p)
}

type params struct {
	fx.In
	Repository  datasetapi.Repository
	DatasetRepo dataset.Repository
	ImageRouter pgin.Router `name:"ImageService"`
}

func provideService(p params) pgin.Router {
	return datasetapi.NewService(p.Repository, p.DatasetRepo, p.ImageRouter)
}
