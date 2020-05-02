package gin

import "github.com/gin-gonic/gin"

type IEngine interface {
	Register(router gin.IRouter)
}
