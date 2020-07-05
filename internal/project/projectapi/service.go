package projectapi

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"

	"github.com/nkhang/pluto/pkg/errors"
	pgin "github.com/nkhang/pluto/pkg/gin"
	"github.com/nkhang/pluto/pkg/ginwrapper"
)

type service struct {
	repository     Repository
	labelService   pgin.IEngine
	datasetService pgin.IEngine
}

const (
	FieldProjectID = "project_id"
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
	router.POST("", ginwrapper.Wrap(s.create))
	router.GET("/:"+FieldProjectID, ginwrapper.Wrap(s.get))
	router.POST("/:"+FieldProjectID+"/perm", ginwrapper.Wrap(s.createPerm))
}

func (s *service) getAll(c *gin.Context) ginwrapper.Response {
	var req GetProjectParam
	if err := c.ShouldBindQuery(&req); err != nil {
		return ginwrapper.Response{
			Error: errors.BadRequest.NewWithMessage("error binding params"),
		}
	}
	responses, total, err := s.repository.GetList(req)
	if err != nil {
		return ginwrapper.Response{
			Error: err,
		}
	}
	return ginwrapper.Response{
		Error: errors.Success.NewWithMessage("Success"),
		Data: GetProjectResponse{
			Total:    total,
			Projects: responses,
		},
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

func (s *service) create(c *gin.Context) ginwrapper.Response {
	var req CreateProjectParams
	if err := c.ShouldBind(&req); err != nil {
		return ginwrapper.Response{
			Error: errors.BadRequest.Wrap(err, "cannot bind request params"),
		}
	}
	err := s.repository.Create(req)
	if err != nil {
		return ginwrapper.Response{
			Error: err,
		}
	}
	return ginwrapper.Response{
		Error: errors.Success.NewWithMessage("success"),
	}
}

func (s *service) createPerm(c *gin.Context) ginwrapper.Response {
	var req CreatePermParams
	idStr := c.Param(FieldProjectID)
	pID, err := cast.ToUint64E(idStr)
	if err != nil {
		return ginwrapper.Response{
			Error: errors.BadRequest.Wrap(err, "cannot get project id"),
		}
	}
	if err := c.ShouldBind(&req); err != nil {
		return ginwrapper.Response{
			Error: errors.BadRequest.NewWithMessage("error binding request params"),
		}
	}
	req.ProjectID = pID
	err = s.repository.CreatePerm(req)
	if err != nil {
		return ginwrapper.Response{
			Error: err,
		}
	}
	return ginwrapper.Response{
		Error: errors.Success.NewWithMessage("success"),
	}
}
