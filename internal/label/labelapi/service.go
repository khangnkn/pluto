package labelapi

import (
	"github.com/gin-gonic/gin"

	"github.com/nkhang/pluto/pkg/errors"
	"github.com/nkhang/pluto/pkg/ginwrapper"
	"github.com/nkhang/pluto/pkg/logger"
)

type service struct {
	repository Repository
}

func NewService(r Repository) *service {
	return &service{
		repository: r,
	}
}

func (s *service) Register(router gin.IRouter) {
	router.GET("", ginwrapper.Wrap(s.getByProjectID))
	router.POST("", ginwrapper.Wrap(s.create))
}

func (s *service) getByProjectID(c *gin.Context) ginwrapper.Response {
	req := LabelRequest{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		logger.Error(err)
		return ginwrapper.Response{
			Error: errors.BadRequest.NewWithMessage("error binding request"),
		}
	}
	responses, err := s.repository.GetByProject(req.ProjectID)
	if err != nil {
		return ginwrapper.Response{
			Error: err,
		}
	}
	return ginwrapper.Response{
		Data:  responses,
		Error: errors.Success.NewWithMessage("Success"),
	}
}

func (s *service) create(c *gin.Context) ginwrapper.Response {
	var req CreateLabelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("error binding request", err)
		return ginwrapper.Response{
			Error: errors.BadRequest.NewWithMessageF("error binding request. error %v", err),
		}
	}
	if err := s.repository.CreateLabel(req); err != nil {
		return ginwrapper.Response{
			Error: err,
		}
	}
	return ginwrapper.Response{
		Error: errors.Success.NewWithMessage("success"),
	}
}
