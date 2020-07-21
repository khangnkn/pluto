package taskfx

import (
	"github.com/jinzhu/gorm"
	"github.com/nkhang/pluto/internal/dataset/datasetapi"
	"github.com/nkhang/pluto/internal/image"
	"github.com/nkhang/pluto/internal/project/projectapi"
	"github.com/nkhang/pluto/internal/task"
	"github.com/nkhang/pluto/internal/task/taskapi"
	"github.com/nkhang/pluto/pkg/annotation"
	"github.com/nkhang/pluto/pkg/cache"
	pgin "github.com/nkhang/pluto/pkg/pgin"
)

func provideTaskDBRepo(db *gorm.DB) task.DBRepository {
	return task.NewDBRepository(db)
}

func provideTaskRepo(dbRepo task.DBRepository, cacheRepo cache.Cache) task.Repository {
	return task.NewRepository(dbRepo, cacheRepo)
}

func provideAPIRepo(r task.Repository, ir image.Repository, datasetRepo datasetapi.Repository, projectRepo projectapi.Repository, annotationService annotation.Service) taskapi.Repository {
	return taskapi.NewRepository(r, ir, datasetRepo, projectRepo, annotationService)
}

func provideService(r taskapi.Repository, tr task.Repository) (pgin.Router, pgin.StandaloneRouter, *taskapi.Service) {
	service := taskapi.NewService(r, tr)
	return service, service, service
}
