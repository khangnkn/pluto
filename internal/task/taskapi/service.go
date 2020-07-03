package taskapi

import (
	"github.com/gin-gonic/gin"
	"github.com/nkhang/pluto/pkg/errors"
	"github.com/nkhang/pluto/pkg/ginwrapper"
)

type service struct {
	repository Repository
}

func NewService(r Repository) *service {
	return &service{repository: r}
}

func (s *service) Register(router gin.IRouter) {
	router.POST("/", ginwrapper.Wrap(s.createTask))
}

func (s *service) createTask(c *gin.Context) ginwrapper.Response {
	var req CreateTaskRequest
	if err := c.ShouldBind(&req); err != nil {
		return ginwrapper.Response{
			Error: errors.BadRequest.NewWithMessage("error binding create task request"),
		}
	}
	err := s.repository.CreateTask(req)
	if err != nil {
		return ginwrapper.Response{
			Error: err,
		}
	}
	return ginwrapper.Response{
		Error: errors.Success.NewWithMessage("success"),
	}
}
