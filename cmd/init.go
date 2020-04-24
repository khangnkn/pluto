package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/nkhang/pluto/internal/ping"
	"github.com/nkhang/pluto/internal/toolapi"
	"github.com/nkhang/pluto/pkg/logger"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

type params struct {
	fx.In

	Router         *gin.Engine
	ToolRepository toolapi.Repository
}

func initializer(l fx.Lifecycle, p params) {
	g := p.Router.Group("/pluto/api/v1/tools")
	toolapi.NewService(p.ToolRepository).Register(g)
	l.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				port := viper.GetInt("service.port")
				s := ping.NewService()

				s.Register(g)
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
