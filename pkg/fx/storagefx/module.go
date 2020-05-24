package storagefx

import "go.uber.org/fx"

var Module = fx.Provide(provideObjectStorage)