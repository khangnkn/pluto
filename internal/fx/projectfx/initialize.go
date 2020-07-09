package projectfx

import (
	"github.com/jinzhu/gorm"
	"github.com/nkhang/pluto/internal/workspace/workspaceapi"
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

func provideRepository(r project.DBRepository, c cache.Cache) project.Repository {
	return project.NewRepository(r, c)
}

func provideAPIRepository(r project.Repository, dr dataset.Repository, wr workspaceapi.Repository) projectapi.Repository {
	return projectapi.NewRepository(r, dr, wr)
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
