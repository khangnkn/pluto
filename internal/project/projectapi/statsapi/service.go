package statsapi

import (
	"github.com/gin-gonic/gin"
	"github.com/nkhang/pluto/internal/project/projectapi"
	"github.com/nkhang/pluto/pkg/errors"
	"github.com/nkhang/pluto/pkg/ginwrapper"
	"github.com/nkhang/pluto/pkg/util/idextractor"
)

type service struct {
	repository Repository
}

func NewService(r Repository) *service {
	return &service{repository: r}
}

func (s *service) Register(router gin.IRouter) {
	router.GET("/images", ginwrapper.Wrap(s.getImageStats))
	router.GET("/overall", ginwrapper.Wrap(s.getTaskStats))
	router.GET("/members", ginwrapper.Wrap(s.getMemberStats))
	router.GET("/labels", ginwrapper.Wrap(s.getStatsLabel))
}

func (s *service) getImageStats(c *gin.Context) ginwrapper.Response {
	var req GetDatasetStatsRequest
	projectID, err := idextractor.ExtractUint64Param(c, projectapi.FieldProjectID)
	if err != nil {
		return ginwrapper.Response{Error: err}
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		return ginwrapper.Response{
			Error: errors.BadRequest.NewWithMessage("error binding params"),
		}
	}
	stats, err := s.repository.BuildReport(projectID, req.DatasetID)
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

func (s *service) getTaskStats(c *gin.Context) ginwrapper.Response {
	projectID := uint64(c.GetInt64(projectapi.FieldProjectID))
	stats, err := s.repository.BuildTaskReport(projectID)
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

func (s *service) getMemberStats(c *gin.Context) ginwrapper.Response {
	projectID := uint64(c.GetInt64(projectapi.FieldProjectID))
	resp, err := s.repository.BuildMemberReport(projectID)
	if err != nil {
		return ginwrapper.Response{
			Error: err,
		}
	}
	return ginwrapper.Response{
		Error: errors.Success.NewWithMessage("success"),
		Data:  resp,
	}
}

func (s *service) getStatsLabel(c *gin.Context) ginwrapper.Response {
	projectID := uint64(c.GetInt64(projectapi.FieldProjectID))
	var req GetLabelStatsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		return ginwrapper.Response{
			Error: errors.BadRequest.NewWithMessage("error binding params"),
		}
	}
	stats, err := s.repository.BuildLabelReport(projectID, req.LabelID)
	if err != nil {
		return ginwrapper.Response{Error: err}
	}
	return ginwrapper.Response{
		Error: errors.Success.NewWithMessage("success"),
		Data:  stats,
	}
}
