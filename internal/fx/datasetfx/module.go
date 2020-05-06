package datasetfx

import "go.uber.org/fx"

var Module = fx.Provide(
	provideRepository,
	fx.Annotated{
		Name:   "DatasetService",
		Target: provideService,
	})
