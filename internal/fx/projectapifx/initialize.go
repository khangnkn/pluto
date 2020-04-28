package projectapifx

import (
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"

	"github.com/nkhang/pluto/internal/project"
	"github.com/nkhang/pluto/internal/projectapi"
	"github.com/nkhang/pluto/pkg/cache"
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
