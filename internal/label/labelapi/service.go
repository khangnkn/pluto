package labelapi

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"

	"github.com/nkhang/pluto/internal/project/projectapi"
	"github.com/nkhang/pluto/pkg/errors"
	"github.com/nkhang/pluto/pkg/ginwrapper"
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
	router.GET("/", ginwrapper.Wrap(s.getByProjectID))
}

func (s *service) getByProjectID(c *gin.Context) ginwrapper.Response {
	pIDStr := c.Param(projectapi.FieldProjectID)
	pID, err := cast.ToUint64E(pIDStr)
	if err != nil {
		return ginwrapper.Response{
			Error: errors.BadRequest.Wrap(err, "cannot get project ID"),
		}
	}
	responses, err := s.repository.GetByProject(pID)
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
