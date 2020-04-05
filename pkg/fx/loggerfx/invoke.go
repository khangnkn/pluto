package loggerfx

import (
	"go.uber.org/fx"
)

var Invoke = fx.Invoke(initializer)