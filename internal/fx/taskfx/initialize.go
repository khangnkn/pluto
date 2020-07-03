package taskfx

import (
	"github.com/jinzhu/gorm"
	"github.com/nkhang/pluto/internal/image"
	"github.com/nkhang/pluto/internal/task"
	"github.com/nkhang/pluto/internal/task/taskapi"
	"github.com/nkhang/pluto/pkg/cache"
	pgin "github.com/nkhang/pluto/pkg/gin"
)

func provideTaskDBRepo(db *gorm.DB) task.DBRepository {
	return task.NewDBRepository(db)
}

func provideTaskRepo(dbRepo task.DBRepository, cacheRepo cache.Cache) task.Repository {
	return task.NewRepository(dbRepo, cacheRepo)
}

func provideAPIRepo(r task.Repository, ir image.Repository) taskapi.Repository {
	return taskapi.NewRepository(r, ir)
}

func provideService(r taskapi.Repository) pgin.IEngine {
	return taskapi.NewService(r)
}
