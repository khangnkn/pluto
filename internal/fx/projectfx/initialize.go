package projectfx

import (
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"go.uber.org/fx"

	"github.com/nkhang/pluto/internal/project"
	"github.com/nkhang/pluto/internal/projectapi"
	"github.com/nkhang/pluto/pkg/cache"
	"github.com/nkhang/pluto/pkg/gin"
)

func provideProjectDBRepository(db *gorm.DB) project.DiskRepository {
	return project.NewDiskRepository(db)
}

func provideRepository(r project.DiskRepository, client redis.UniversalClient) project.Repository {
	c := cache.New(client)
	return project.NewRepository(r, c)
}

func provideAPIRepository(r project.Repository) projectapi.Repository {
	return projectapi.NewRepository(r)
}

type params struct {
	fx.In

	Repository   projectapi.Repository
	LabelService gin.IEngine `name:"LabelService"`
}

func provideService(p params) gin.IEngine {
	return projectapi.NewService(p.Repository, p.LabelService)
}
