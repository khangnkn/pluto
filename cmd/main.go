package main

import (
	"github.com/nkhang/pluto/pkg/ginfx"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		ginfx.Module,
		fx.Invoke(initializer)
	).Run()
}
