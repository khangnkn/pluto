package pgin

import "github.com/gin-gonic/gin"

type StandaloneRouter interface {
	RegisterStandalone(router gin.IRouter)
}

type Router interface {
	Register(router gin.IRouter)
}
