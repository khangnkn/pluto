package workspaceapi

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"

	"github.com/nkhang/pluto/pkg/errors"
	pgin "github.com/nkhang/pluto/pkg/gin"
	"github.com/nkhang/pluto/pkg/ginwrapper"
)

const (
	FieldWorkspaceID = "workspaceId"
)

type service struct {
	repository     Repository
	projectService pgin.IEngine
}

func NewService(r Repository, pS pgin.IEngine) *service {
	return &service{
		repository:     r,
		projectService: pS,
	}
}

func (s *service) Register(router gin.IRouter) {
	router.GET("/:"+FieldWorkspaceID, ginwrapper.Wrap(s.get))
	projectRouter := router.Group("/:" + FieldWorkspaceID + "/projects")
	s.projectService.Register(projectRouter)
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
		Error: errors.Success.NewWithMessage("Successfully"),
		Data:  w,
	}
}
