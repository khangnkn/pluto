package ping

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type service struct {}

func NewService() *service {
	return &service{}
}

func (s *service) Register(g gin.IRoutes) {
	g.GET("/ping", ping)
}

func ping(c *gin.Context)  {
	c.String(http.StatusOK, "pong!")
}