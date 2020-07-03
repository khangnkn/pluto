package taskfx

import "go.uber.org/fx"

var Module = fx.Provide(
	provideTaskDBRepo,
	provideTaskRepo,
	provideAPIRepo,
	fx.Annotated{
		Name:   "TaskService",
		Target: provideService,
	})
