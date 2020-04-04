package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nkhang/pluto/internal/ping"
	"go.uber.org/fx"
)

type Param struct {
	fx.In

	router *gin.Engine
	server http.Server
}

func initializer(l fx.Lifecycle, p Param) {
	l.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				s := ping.NewService()
				s.Register(p.router)
				p.router.Run(":8080")
				return nil
			},
			OnStop: func(ctx context.Context) error {
				log.Println("Stop Server")
				return nil
			},
		},
	)
}
