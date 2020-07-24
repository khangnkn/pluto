package taskfx

import (
	"github.com/jinzhu/gorm"
	"github.com/nkhang/pluto/internal/dataset/datasetapi"
	"github.com/nkhang/pluto/internal/image"
	"github.com/nkhang/pluto/internal/project/projectapi"
	"github.com/nkhang/pluto/internal/task"
	"github.com/nkhang/pluto/internal/task/taskapi"
	"github.com/nkhang/pluto/internal/task/taskapi/statsapi"
	"github.com/nkhang/pluto/pkg/annotation"
	"github.com/nkhang/pluto/pkg/cache"
	pgin "github.com/nkhang/pluto/pkg/pgin"
	"go.uber.org/fx"
)

func provideTaskDBRepo(db *gorm.DB) task.DBRepository {
	return task.NewDBRepository(db)
}

func provideTaskRepo(dbRepo task.DBRepository, cacheRepo cache.Cache) task.Repository {
	return task.NewRepository(dbRepo, cacheRepo)
}

func provideTaskStatsRepo(r task.Repository) statsapi.Repository {
	return statsapi.NewRepository(r)
}

func provideTaskStatsService(r statsapi.Repository) pgin.Router {
	return statsapi.NewService(r)
}

func provideAPIRepo(r task.Repository, ir image.Repository, datasetRepo datasetapi.Repository, projectRepo projectapi.Repository, annotationService annotation.Service) taskapi.Repository {
	return taskapi.NewRepository(r, ir, datasetRepo, projectRepo, annotationService)
}

type params struct {
	fx.In
	Repo        task.Repository
	APIRepo     taskapi.Repository
	StatsRouter pgin.Router `name:"TaskStatsService"`
}

func provideService(p params) (pgin.Router, pgin.StandaloneRouter, *taskapi.Service) {
	service := taskapi.NewService(p.APIRepo, p.Repo, p.StatsRouter)
	return service, service, service
}
