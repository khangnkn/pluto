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
