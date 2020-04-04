package ping

import "github.com/gin-gonic/gin"

func initializer(g gin.IRoutes) {
	s := NewService()
	s.Register(g)
}