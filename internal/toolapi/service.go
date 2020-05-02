package toolapi

import (
	"github.com/gin-gonic/gin"
	"github.com/nkhang/pluto/pkg/errors"
	"github.com/nkhang/pluto/pkg/ginwrapper"
)

type service struct {
	repo Repository
}

func NewService(r Repository) *service {
	return &service{
		repo: r,
	}
}

func (s *service) Register(router gin.IRouter) {
	router.GET("/", ginwrapper.Wrap(s.allTools))
}

func (s *service) allTools(c *gin.Context) ginwrapper.Response {
	tools, err := s.repo.GetAll()
	if err != nil {
		return ginwrapper.Response{
			Error: err,
		}
	}
	return ginwrapper.Response{
		Error: errors.Success.NewWithMessage("Successfully transaction"),
		Data:  tools,
	}
}
