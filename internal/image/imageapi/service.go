package imageapi

import (
	"github.com/gin-gonic/gin"
	"github.com/nkhang/pluto/internal/dataset/datasetapi"
	"github.com/nkhang/pluto/pkg/errors"
	"github.com/nkhang/pluto/pkg/ginwrapper"
	"github.com/spf13/cast"
)

type service struct {
	repository Repository
}

func NewService(r Repository) *service {
	return &service{repository: r}
}

const fieldImageID = "imageId"

func (s *service) Register(router gin.IRouter) {
	router.GET("", ginwrapper.Wrap(s.getByDataset))
	router.POST("", ginwrapper.Wrap(s.uploadByDataset))
	router.GET("/:"+fieldImageID, ginwrapper.Wrap(s.get))
}

func (s *service) RegisterStandalone(router gin.IRouter) {
	router.GET("/:"+fieldImageID, ginwrapper.Wrap(s.get))
}

func (s *service) get(c *gin.Context) ginwrapper.Response {
	idStr := c.Param(fieldImageID)
	id, err := cast.ToUint64E(idStr)
	if err != nil {
		return ginwrapper.Response{
			Error: errors.BadRequest.Wrap(err, "error binding params"),
		}
	}
	req := GetImageRequest{ID: id}
	response, err := s.repository.GetImage(req)
	if err != nil {
		return ginwrapper.Response{
			Error: err,
		}
	}
	return ginwrapper.Response{
		Error: errors.Success.NewWithMessage("success"),
		Data:  response,
	}
}

func (s *service) getByDataset(c *gin.Context) ginwrapper.Response {
	datasetID := uint64(c.GetInt64(datasetapi.FieldDatasetID))
	var q ImageRequestQuery
	err := c.ShouldBindQuery(&q)
	if err != nil {
		return ginwrapper.Response{
			Error: errors.BadRequest.Wrap(err, "fail to bind request"),
		}
	}
	responses, err := s.repository.GetByDatasetID(datasetID, q.Offset, q.Limit)
	if err != nil {
		return ginwrapper.Response{
			Error: err,
		}
	}
	return ginwrapper.Response{
		Error: errors.Success.NewWithMessage("success"),
		Data:  responses,
	}
}

func (s *service) uploadByDataset(c *gin.Context) ginwrapper.Response {
	datasetID := uint64(c.GetInt64(datasetapi.FieldDatasetID))
	var req UploadRequest
	err := c.ShouldBind(&req)
	if err != nil {
		return ginwrapper.Response{
			Error: errors.BadRequest.NewWithMessage("error binding request"),
		}
	}
	d, errs := s.repository.UploadRequest(datasetID, req.FileHeader)
	if len(errs) != 0 {
		return ginwrapper.Response{
			Error: errors.ImageErrorCreating.Wrap(err, "error reading file"),
			Data:  errs,
		}
	}
	return ginwrapper.Response{
		Error: errors.Success.NewWithMessage("success"),
		Data:  d,
	}
}
