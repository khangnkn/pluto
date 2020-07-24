package taskfx

import "go.uber.org/fx"

var Module = fx.Provide(
	provideTaskDBRepo,
	provideTaskRepo,
	provideAPIRepo,
	provideTaskStatsRepo,
	fx.Annotated{
		Name:   "TaskService",
		Target: provideService,
	},
	fx.Annotated{
		Name:   "TaskStatsService",
		Target: provideTaskStatsService,
	})
