package projectfx

import (
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"go.uber.org/fx"

	"github.com/nkhang/pluto/internal/dataset"
	"github.com/nkhang/pluto/internal/project"
	"github.com/nkhang/pluto/internal/project/projectapi"
	"github.com/nkhang/pluto/pkg/cache"
	"github.com/nkhang/pluto/pkg/gin"
)

func provideProjectDBRepository(db *gorm.DB) project.DBRepository {
	return project.NewDiskRepository(db)
}

func provideRepository(r project.DBRepository, client redis.UniversalClient) project.Repository {
	c := cache.New(client)
	return project.NewRepository(r, c)
}

func provideAPIRepository(r project.Repository, dr dataset.Repository) projectapi.Repository {
	return projectapi.NewRepository(r, dr)
}

type params struct {
	fx.In

	Repository     projectapi.Repository
	LabelService   gin.IEngine `name:"LabelService"`
	DatasetService gin.IEngine `name:"DatasetService"`
}

func provideService(p params) gin.IEngine {
	return projectapi.NewService(p.Repository, p.LabelService, p.DatasetService)
}
