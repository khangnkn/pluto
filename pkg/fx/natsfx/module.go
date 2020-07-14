package natsfx

import "go.uber.org/fx"

var Module = fx.Provide(provideNATSClient)
