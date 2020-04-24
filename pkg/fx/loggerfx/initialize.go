package loggerfx

import (
	"github.com/nkhang/pluto/pkg/logger"
	"github.com/spf13/viper"
)

func initializer() {
	prod := viper.GetBool("service.production")
	logger.Initlialize(prod)
}
