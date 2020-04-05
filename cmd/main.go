package main

import (
	"github.com/nkhang/pluto/pkg/fx/configfx"
	"github.com/nkhang/pluto/pkg/fx/ginfx"
	"github.com/nkhang/pluto/pkg/fx/loggerfx"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		configfx.Initialize("pluto"),
		loggerfx.Invoke,
		ginfx.Module,
		fx.Invoke(initializer),
	).Run()
}
