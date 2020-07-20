package statsapi

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
	router.GET("/dataset", ginwrapper.Wrap(s.getImageStats))
}

func (s *service) getImageStats(c *gin.Context) ginwrapper.Response {
	var req GetDatasetStatsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		return ginwrapper.Response{
			Error: errors.BadRequest.NewWithMessage("error binding params"),
		}
	}
	stats, err := s.repository.BuildReport(req.DatasetID)
	if err != nil {
		return ginwrapper.Response{
			Error: err,
		}
	}
	return ginwrapper.Response{
		Error: errors.Success.NewWithMessage("success"),
		Data:  stats,
	}
}
