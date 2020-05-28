package projectapi

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"

	"github.com/nkhang/pluto/pkg/errors"
	pgin "github.com/nkhang/pluto/pkg/gin"
	"github.com/nkhang/pluto/pkg/ginwrapper"
	"github.com/nkhang/pluto/pkg/logger"
)

type service struct {
	repository     Repository
	labelService   pgin.IEngine
	datasetService pgin.IEngine
}

const (
	FieldProjectID = "projectId"
)

func NewService(r Repository, labelService, datasetService pgin.IEngine) *service {
	return &service{
		repository:     r,
		labelService:   labelService,
		datasetService: datasetService,
	}
}

func (s *service) Register(router gin.IRouter) {
	router.GET("", ginwrapper.Wrap(s.getAll))
	router.GET("/:"+FieldProjectID, ginwrapper.Wrap(s.get))
}

func (s *service) getAll(c *gin.Context) ginwrapper.Response {
	var req GetProjectParam
	if err := c.ShouldBindQuery(&req); err != nil {
		return ginwrapper.Response{
			Error: errors.BadRequest.NewWithMessage("error binding params"),
		}
	}
	responses, err := s.repository.GetByWorkspaceID(req.WorkspaceID)
	if err != nil {
		logger.Error(err)
		return ginwrapper.Response{
			Error: err,
		}
	}
	return ginwrapper.Response{
		Error: errors.Success.NewWithMessage("Success"),
		Data:  responses,
	}
}

func (s *service) get(c *gin.Context) ginwrapper.Response {
	idStr := c.Param(FieldProjectID)
	pID, err := cast.ToUint64E(idStr)
	if err != nil {
		return ginwrapper.Response{
			Error: errors.BadRequest.Wrap(err, "cannot get project id"),
		}
	}
	if pID <= 0 {
		return ginwrapper.Response{
			Error: errors.BadRequest.NewWithMessage("id must greater than 0"),
		}
	}
	p, err := s.repository.GetByID(pID)
	if err != nil {
		return ginwrapper.Response{
			Error: err,
		}
	}
	return ginwrapper.Response{
		Error: errors.Success.NewWithMessage("Success"),
		Data:  p,
	}
}

func getWorkspaceParams(c *gin.Context) (uint64, error) {
	wsIdStr := c.Param(FieldProjectID)
	return cast.ToUint64E(wsIdStr)
}
