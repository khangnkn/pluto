package datasetapi

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"

	"github.com/nkhang/pluto/internal/project/projectapi"
	"github.com/nkhang/pluto/pkg/errors"
	"github.com/nkhang/pluto/pkg/ginwrapper"
)

const (
	fieldDatasetID = "datasetId"
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
	router.GET("/:"+fieldDatasetID, ginwrapper.Wrap(s.getByID))
}

func (s *service) getByID(c *gin.Context) ginwrapper.Response {
	dIDStr := c.Param(fieldDatasetID)
	dID, err := cast.ToUint64E(dIDStr)
	if err != nil {
		return ginwrapper.Response{
			Error: errors.BadRequest.Wrap(err, "cannot get dataset id"),
		}
	}
	dataset, err := s.repository.GetByID(dID)
	if err != nil {
		return ginwrapper.Response{
			Error: err,
		}
	}
	return ginwrapper.Response{
		Error: errors.Success.NewWithMessage("Success"),
		Data:  dataset,
	}
}

func (s *service) getByProjectID(c *gin.Context) ginwrapper.Response {
	pIDStr := c.Param(projectapi.FieldProjectID)
	pID, err := cast.ToUint64E(pIDStr)
	if err != nil {
		return ginwrapper.Response{
			Error: errors.BadRequest.Wrap(err, "cannot get dataset id"),
		}
	}
	datasets, err := s.repository.GetByProjectID(pID)
	if err != nil {
		return ginwrapper.Response{
			Error: err,
		}
	}
	return ginwrapper.Response{
		Error: errors.Success.NewWithMessage("Success"),
		Data:  datasets,
	}
}
