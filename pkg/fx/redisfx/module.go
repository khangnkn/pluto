package redisfx

import "go.uber.org/fx"

var Module = fx.Provide(provideRedisClient, provideCacheRepository)
