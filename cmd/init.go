package main

import (
	"context"
	"fmt"
	"log"

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
	pgin "github.com/nkhang/pluto/pkg/gin"
	"github.com/nkhang/pluto/pkg/logger"
)

type params struct {
	fx.In

	GormDB           *gorm.DB
	Router           *gin.Engine
	ToolRepository   toolapi.Repository
	WorkspaceService pgin.IEngine `name:"WorkspaceService"`
	ProjectService   pgin.IEngine `name:"ProjectService"`
	DatasetService   pgin.IEngine `name:"DatasetService"`
	LabelService     pgin.IEngine `name:"LabelService"`
	ToolService      pgin.IEngine `name:"ToolService"`
	ImageService     pgin.IEngine `name:"ImageService"`
}

func initializer(l fx.Lifecycle, p params) {
	migrate(p.GormDB)
	router := p.Router.Group("/pluto/api/v1")
	p.ToolService.Register(router.Group("/tools"))
	p.LabelService.Register(router.Group("/labels"))
	p.ProjectService.Register(router.Group("/projects"))
	p.DatasetService.Register(router.Group("/datasets"))
	p.WorkspaceService.Register(router.Group("/workspaces"))
	p.ImageService.Register(router.Group("/images"))
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
}
