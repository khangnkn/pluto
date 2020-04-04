package ginfx

import (
	"github.com/gin-gonic/gin"
)

func initializer() (*gin.Engine, gin.IRouter) {
	e := gin.Default()
	return e, e
}
