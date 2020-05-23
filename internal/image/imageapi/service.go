package imageapi

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
	router.GET("", ginwrapper.Wrap(s.getByDataset))
	router.POST("", ginwrapper.Wrap(s.uploadByDataset))
}

func (s *service) getByDataset(c *gin.Context) ginwrapper.Response {
	var q ImageRequestQuery
	err := c.ShouldBindQuery(&q)
	if err != nil {
		return ginwrapper.Response{
			Error: errors.BadRequest.Wrap(err, "fail to bind request"),
		}
	}
	responses, err := s.repository.GetByDatasetID(q.DatasetID, q.Offset, q.Limit)
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
	var req UpdloadRequest
	file, err := c.FormFile("file")
	if err != nil {
		return ginwrapper.Response{
			Error: errors.BadRequest.NewWithMessage("error getting image"),
		}
	}
	err = c.ShouldBind(&req)
	if err != nil {
		return ginwrapper.Response{
			Error: errors.BadRequest.NewWithMessage("error binding request"),
		}
	}
	err = s.repository.UploadRequest(1, file)
	if err != nil {
		return ginwrapper.Response{
			Error: errors.DatasetQueryError.Wrap(err, "error reading file"),
		}
	}
	return ginwrapper.Response{
		Error: errors.Success.NewWithMessage("success"),
	}
}
