package statsapi

import (
	"github.com/gin-gonic/gin"
	"github.com/nkhang/pluto/internal/task/taskapi"
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
	router.GET("/stats", ginwrapper.Wrap(s.getTaskStats))
}

func (s *service) getTaskStats(c *gin.Context) ginwrapper.Response {
	taskID, err := idextractor.ExtractUint64Param(c, taskapi.FieldTaskID)
	if err != nil {
		return ginwrapper.Response{
			Error: err,
		}
	}
	stats, err := s.repository.Stats(taskID)
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
