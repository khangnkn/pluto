package workspacefx

import (
	"github.com/jinzhu/gorm"
	"github.com/nkhang/pluto/internal/project"
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

func provideWorkspaceAPIRepository(workspaceRepo workspace.Repository, projectRepo project.Repository) workspaceapi.Repository {
	return workspaceapi.NewRepository(workspaceRepo, projectRepo)
}

type params struct {
	fx.In
	Repository workspaceapi.Repository
}

func provideWorkspaceService(p params) gin.IEngine {
	return workspaceapi.NewService(p.Repository)
}
