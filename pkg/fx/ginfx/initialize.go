package ginfx

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func initializer() (*gin.Engine, gin.IRouter) {
	prod := viper.GetBool("service.production")
	if prod {
		gin.SetMode(gin.ReleaseMode)
	}
	e := gin.Default()
	conf := cors.DefaultConfig()
	conf.AllowOrigins = append(conf.AllowOrigins, "http://localhost:3000")
	conf.AllowCredentials = true
	conf.AllowFiles = true
	conf.AddAllowHeaders("Authorization")
	e.Use(cors.New(conf))
	return e, e
}
