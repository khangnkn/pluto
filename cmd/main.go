package main

import (
	"github.com/nkhang/pluto/internal/fx/labelfx"
	"github.com/nkhang/pluto/internal/fx/projectfx"
	"github.com/nkhang/pluto/internal/fx/toolfx"
	"github.com/nkhang/pluto/internal/fx/workspacefx"
	"github.com/nkhang/pluto/pkg/fx/configfx"
	"github.com/nkhang/pluto/pkg/fx/dbfx"
	"github.com/nkhang/pluto/pkg/fx/ginfx"
	"github.com/nkhang/pluto/pkg/fx/loggerfx"
	"github.com/nkhang/pluto/pkg/fx/redisfx"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		configfx.Initialize("pluto"),
		loggerfx.Invoke,
		dbfx.Module,
		redisfx.Module,
		toolfx.Module,
		labelfx.Module,
		projectfx.Module,
		workspacefx.Module,
		ginfx.Module,
		fx.Invoke(initializer),
	).Run()
}
