package projectfx

import "go.uber.org/fx"

var Module = fx.Provide(
	provideProjectDBRepository,
	provideRepository,
	provideAPIRepository,
	provideStatsAPIRepo,
	fx.Annotated{
		Name:   "ProjectService",
		Target: provideService,
	},
)
