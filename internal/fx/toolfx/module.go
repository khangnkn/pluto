package toolfx

import "go.uber.org/fx"

var Module = fx.Provide(
	provideToolRepository,
	provideToolAPI,
	fx.Annotated{
		Name:   "ToolService",
		Target: provideToolService,
	},
)
