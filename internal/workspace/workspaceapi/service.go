package workspaceapi

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"

	"github.com/nkhang/pluto/pkg/errors"
	"github.com/nkhang/pluto/pkg/ginwrapper"
)

const (
	FieldWorkspaceID = "workspace_id"
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
	router.GET("/", ginwrapper.Wrap(s.getByUserID))
	router.POST("/", ginwrapper.Wrap(s.create))
	router.GET("/:"+FieldWorkspaceID, ginwrapper.Wrap(s.get))
}

func (s *service) get(c *gin.Context) ginwrapper.Response {
	idStr := c.Param(FieldWorkspaceID)
	id, err := cast.ToUint64E(idStr)
	if err != nil {
		return ginwrapper.Response{
			Error: errors.BadRequest.Wrap(err, "invalid id"),
		}
	}
	w, err := s.repository.GetByID(id)
	if err != nil {
		return ginwrapper.Response{
			Error: err,
		}
	}
	return ginwrapper.Response{
		Error: errors.Success.NewWithMessage("success"),
		Data:  w,
	}
}

func (s *service) getByUserID(c *gin.Context) ginwrapper.Response {
	var req GetByUserIDRequest
	err := c.ShouldBindQuery(&req)
	if err != nil {
		return ginwrapper.Response{
			Error: errors.BadRequest.NewWithMessage("error binding user_id"),
		}
	}
	workspaces, err := s.repository.GetByUserID(req)
	if err != nil {
		return ginwrapper.Response{
			Error: err,
		}
	}
	return ginwrapper.Response{
		Error: errors.Success.NewWithMessage("success"),
		Data:  workspaces,
	}
}

func (s *service) create(c *gin.Context) ginwrapper.Response {
	var req CreateWorkspaceRequest
	if err := c.ShouldBind(&req); err != nil {
		return ginwrapper.Response{
			Error: errors.BadRequest.Wrap(err, "cannot bind create workspace body"),
		}
	}
	response, err := s.repository.CreateWorkspace(req)
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
