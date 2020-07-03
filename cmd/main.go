package main

import (
	"github.com/nkhang/pluto/internal/fx/taskfx"
	"go.uber.org/fx"

	"github.com/nkhang/pluto/internal/fx/datasetfx"
	"github.com/nkhang/pluto/internal/fx/imagefx"
	"github.com/nkhang/pluto/internal/fx/labelfx"
	"github.com/nkhang/pluto/internal/fx/projectfx"
	"github.com/nkhang/pluto/internal/fx/toolfx"
	"github.com/nkhang/pluto/internal/fx/workspacefx"
	"github.com/nkhang/pluto/pkg/fx/configfx"
	"github.com/nkhang/pluto/pkg/fx/dbfx"
	"github.com/nkhang/pluto/pkg/fx/ginfx"
	"github.com/nkhang/pluto/pkg/fx/redisfx"
	"github.com/nkhang/pluto/pkg/fx/storagefx"
)

func main() {
	fx.New(
		configfx.Initialize("pluto"),
		dbfx.Module,
		redisfx.Module,
		toolfx.Module,
		labelfx.Module,
		taskfx.Module,
		imagefx.Module,
		datasetfx.Module,
		projectfx.Module,
		workspacefx.Module,
		storagefx.Module,
		ginfx.Module,
		fx.Invoke(initializer),
	).Run()
}
