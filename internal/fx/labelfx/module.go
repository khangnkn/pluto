package labelfx

import "go.uber.org/fx"

var Module = fx.Provide(
	provideRepository,
	fx.Annotated{
		Name:   "LabelService",
		Target: provideService,
	})
