package datasetapi

import (
	"github.com/gin-gonic/gin"
	"github.com/nkhang/pluto/pkg/util/idextractor"
	"github.com/spf13/cast"

	"github.com/nkhang/pluto/pkg/errors"
	"github.com/nkhang/pluto/pkg/ginwrapper"
	"github.com/nkhang/pluto/pkg/logger"
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

func (s *service) RegisterStandalone(router gin.IRouter) {
	router.GET("", ginwrapper.Wrap(s.getByProjectID))
	router.GET("/:"+fieldDatasetID, ginwrapper.Wrap(s.getByID))
	router.DELETE("/:"+fieldDatasetID, ginwrapper.Wrap(s.del))
	router.POST("/:"+fieldDatasetID+"/clone", ginwrapper.Wrap(s.clone))
	router.POST("", ginwrapper.Wrap(s.create))
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
	var request GetDatasetRequest
	if err := c.ShouldBind(&request); err != nil {
		return ginwrapper.Response{
			Error: errors.BadRequest.Wrap(err, "cannot get dataset id"),
		}
	}
	datasets, err := s.repository.GetByProjectID(request.ProjectID)
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

func (s *service) create(c *gin.Context) ginwrapper.Response {
	var req CreateDatasetRequest
	if err := c.ShouldBind(&req); err != nil {
		return ginwrapper.Response{
			Error: errors.BadRequest.Wrap(err, "cannot bind request"),
		}
	}
	dataset, err := s.repository.CreateDataset(req.Title, req.Description, req.ProjectID)
	if err != nil {
		return ginwrapper.Response{
			Error: err,
		}
	}
	return ginwrapper.Response{
		Error: errors.Success.NewWithMessage("success"),
		Data:  dataset,
	}
}

func (s *service) clone(c *gin.Context) ginwrapper.Response {
	var req CloneDatasetRequest
	if err := c.ShouldBind(&req); err != nil {
		return ginwrapper.Response{
			Error: errors.BadRequest.Wrap(err, "cannot bind request"),
		}
	}
	dIDStr := c.Param(fieldDatasetID)
	dID, err := cast.ToUint64E(dIDStr)
	if err != nil {
		return ginwrapper.Response{
			Error: errors.BadRequest.Wrap(err, "cannot get dataset id"),
		}
	}
	logger.Info("start cloning project")
	cloned, err := s.repository.CloneDataset(req.ProjectID, dID)
	if err != nil {
		return ginwrapper.Response{
			Error: err,
		}
	}
	return ginwrapper.Response{
		Error: errors.Success.NewWithMessage("success"),
		Data:  cloned,
	}
}

func (s *service) del(c *gin.Context) ginwrapper.Response {
	id, err := idextractor.ExtractUint64Param(c, fieldDatasetID)
	if err != nil {
		return ginwrapper.Response{Error: err}
	}
	err = s.repository.Delete(id)
	if err != nil {
		return ginwrapper.Response{Error: err}
	}
	return ginwrapper.Response{
		Error: errors.Success.NewWithMessage("success"),
	}
}
