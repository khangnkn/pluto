package projectfx

import "go.uber.org/fx"

var Module = fx.Provide(
	provideProjectDBRepository,
	provideRepository,
	provideAPIRepository,
	fx.Annotated{
		Name:   "ProjectService",
		Target: provideService,
	},
)
