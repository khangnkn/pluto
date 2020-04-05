package ginfx

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func initializer() (*gin.Engine, gin.IRouter) {
	prod := viper.GetBool("service.production")
	if prod {
		gin.SetMode(gin.ReleaseMode)
	}
	e := gin.Default()
	return e, e
}
