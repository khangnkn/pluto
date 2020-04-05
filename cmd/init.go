package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nkhang/pluto/internal/ping"
	"github.com/nkhang/pluto/pkg/logger"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"log"
)

type Param struct {
	fx.In

	Router *gin.Engine
}

func initializer(l fx.Lifecycle, p Param) {
	l.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				port := viper.GetInt("service.port")
				s := ping.NewService()
				s.Register(p.Router)
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
