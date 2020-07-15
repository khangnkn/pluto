package main

import (
	"context"
	"fmt"
	"log"

	"github.com/nkhang/pluto/internal/task"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"go.uber.org/fx"

	"github.com/nkhang/pluto/internal/dataset"
	"github.com/nkhang/pluto/internal/image"
	"github.com/nkhang/pluto/internal/label"
	"github.com/nkhang/pluto/internal/project"
	"github.com/nkhang/pluto/internal/tool"
	"github.com/nkhang/pluto/internal/tool/toolapi"
	"github.com/nkhang/pluto/internal/workspace"
	"github.com/nkhang/pluto/pkg/logger"
	pgin "github.com/nkhang/pluto/pkg/pgin"
)

type params struct {
	fx.In

	GormDB           *gorm.DB
	Router           *gin.Engine
	ToolRepository   toolapi.Repository
	WorkspaceService pgin.StandaloneRouter `name:"WorkspaceService"`
	ProjectService   pgin.StandaloneRouter `name:"ProjectService"`
	DatasetService   pgin.StandaloneRouter `name:"DatasetService"`
	LabelService     pgin.StandaloneRouter `name:"LabelService"`
	ToolService      pgin.StandaloneRouter `name:"ToolService"`
	ImageService     pgin.StandaloneRouter `name:"ImageService"`
	TaskService      pgin.StandaloneRouter `name:"TaskService"`
}

func initializer(l fx.Lifecycle, p params) {
	migrate(p.GormDB)
	router := p.Router.Group("/pluto/api/v1")
	p.ToolService.RegisterStandalone(router.Group("/tools"))
	p.LabelService.RegisterStandalone(router.Group("/labels"))
	p.ProjectService.RegisterStandalone(router.Group("/projects"))
	p.DatasetService.RegisterStandalone(router.Group("/datasets"))
	p.WorkspaceService.RegisterStandalone(router.Group("/workspaces"))
	p.ImageService.RegisterStandalone(router.Group("/images"))
	p.TaskService.RegisterStandalone(router.Group("/tasks"))
	l.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				port := viper.GetInt("service.port")
				go func() {
					addr := fmt.Sprintf(":%d", port)
					err := p.Router.Run(addr)
					if err != nil {
						logger.Panic(err)
					}
				}()
				logger.Infof("Server is running at port %d", port)
				return nil
			},
			OnStop: func(ctx context.Context) error {
				log.Println("Stop Server")
				return nil
			},
		},
	)
}

func migrate(db *gorm.DB) {
	db.AutoMigrate(&tool.Tool{})
	db.AutoMigrate(&dataset.Dataset{})
	db.AutoMigrate(&label.Label{})
	db.AutoMigrate(&project.Project{})
	db.AutoMigrate(&project.Permission{})
	db.AutoMigrate(&workspace.Workspace{})
	db.AutoMigrate(&workspace.Permission{})
	db.AutoMigrate(&image.Image{})
	db.AutoMigrate(&task.Task{})
	db.AutoMigrate(&task.Detail{})
	db.AutoMigrate(&task.Detail{TaskID: 1})
	db.AutoMigrate(&task.Detail{TaskID: 2})
	db.AutoMigrate(&task.Detail{TaskID: 3})
	db.AutoMigrate(&task.Detail{TaskID: 4})
	db.AutoMigrate(&task.Detail{TaskID: 5})
	db.AutoMigrate(&task.Detail{TaskID: 6})
	db.AutoMigrate(&task.Detail{TaskID: 7})
	db.AutoMigrate(&task.Detail{TaskID: 8})
	db.AutoMigrate(&task.Detail{TaskID: 9})
}
