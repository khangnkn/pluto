package workspacefx

import (
	"github.com/jinzhu/gorm"
	"go.uber.org/fx"

	"github.com/nkhang/pluto/internal/workspace"
	"github.com/nkhang/pluto/internal/workspace/workspaceapi"
	"github.com/nkhang/pluto/pkg/cache"
	"github.com/nkhang/pluto/pkg/gin"
)

func provideWorkspaceDBRepository(db *gorm.DB) workspace.DBRepository {
	return workspace.NewDDBRepository(db)
}

func provideWorkspaceRepository(r workspace.DBRepository, c cache.Cache) workspace.Repository {
	return workspace.NewRepository(r, c)
}

func provideWorkspaceAPIRepository(r workspace.Repository) workspaceapi.Repository {
	return workspaceapi.NewRepository(r)
}

type params struct {
	fx.In
	Repository    workspaceapi.Repository
	ProjectRouter gin.IEngine `name:"ProjectService"`
}

func provideWorkspaceService(p params) gin.IEngine {
	return workspaceapi.NewService(p.Repository, p.ProjectRouter)
}
