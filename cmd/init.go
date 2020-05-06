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
	ToolService      pgin.IEngine `name:"ToolService"`
}

func initializer(l fx.Lifecycle, p params) {
	migrate(p.GormDB)
	router := p.Router.Group("/pluto/api/v1")
	p.ToolService.Register(router.Group("/tools"))
	p.WorkspaceService.Register(router.Group("/workspaces"))
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
	db.AutoMigrate(&workspace.Workspace{})
}
