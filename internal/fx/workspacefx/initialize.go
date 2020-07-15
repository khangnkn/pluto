package workspacefx

import (
	"github.com/jinzhu/gorm"
	"github.com/nkhang/pluto/internal/project"
	"github.com/nkhang/pluto/internal/workspace/workspaceapi/permissionapi"
	"go.uber.org/fx"

	"github.com/nkhang/pluto/internal/workspace"
	"github.com/nkhang/pluto/internal/workspace/workspaceapi"
	"github.com/nkhang/pluto/pkg/cache"
	"github.com/nkhang/pluto/pkg/pgin"
)

func provideWorkspaceDBRepository(db *gorm.DB) workspace.DBRepository {
	return workspace.NewDDBRepository(db)
}

func provideWorkspaceRepository(r workspace.DBRepository, projectRepo project.Repository, c cache.Cache) workspace.Repository {
	return workspace.NewRepository(r, projectRepo, c)
}

func provideWorkspaceAPIRepository(workspaceRepo workspace.Repository, projectRepo project.Repository) workspaceapi.Repository {
	return workspaceapi.NewRepository(workspaceRepo, projectRepo)
}

type params struct {
	fx.In
	Repository    workspaceapi.Repository
	Wr            workspace.Repository
	ProjectRouter pgin.Router `name:"ProjectService"`
}

func provideWorkspaceService(p params) pgin.StandaloneRouter {
	permRepo := permissionapi.NewRepository(p.Wr)
	permRouter := permissionapi.NewService(permRepo)
	return workspaceapi.NewService(p.Repository, p.Wr, p.ProjectRouter, permRouter)
}
