package imagefx

import "go.uber.org/fx"

var Module = fx.Provide(
	provideImageRepository,
	fx.Annotated{
		Name:   "ImageService",
		Target: provideService,
	},
)
