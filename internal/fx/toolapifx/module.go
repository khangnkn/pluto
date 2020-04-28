package toolapifx

import "go.uber.org/fx"

var Module = fx.Provide(
	provideToolRepository,
	provideToolAPI,
)
