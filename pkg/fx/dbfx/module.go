package dbfx

import "go.uber.org/fx"

var Module = fx.Provide(initializer)
