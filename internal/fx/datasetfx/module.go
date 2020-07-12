package datasetfx

import "go.uber.org/fx"

var Module = fx.Provide(
	provideRepository,
	provideAPIRepo,
	fx.Annotated{
		Name:   "DatasetService",
		Target: provideService,
	})
