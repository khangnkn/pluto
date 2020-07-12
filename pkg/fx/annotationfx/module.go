package annotationfx

import "go.uber.org/fx"

var Module = fx.Provide(provideAnnotationService)
