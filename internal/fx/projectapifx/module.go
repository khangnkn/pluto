package projectapifx

import "go.uber.org/fx"

var Module = fx.Provide(
	provideProjectDBRepository,
	provideRepository,
	provideAPIRepository,
)
