package permissionapi

import (
	"github.com/gin-gonic/gin"
	"github.com/nkhang/pluto/internal/workspace/workspaceapi/consts"
	"github.com/nkhang/pluto/pkg/errors"
	"github.com/nkhang/pluto/pkg/ginwrapper"
	"github.com/nkhang/pluto/pkg/logger"
	"github.com/nkhang/pluto/pkg/util/idextractor"
)

type service struct {
	repository Repository
}

func NewService(r Repository) *service {
	return &service{repository: r}
}

const FieldUserID = "userId"

func (s *service) Register(router gin.IRouter) {
	router.GET("", ginwrapper.Wrap(s.get))
	router.POST("", ginwrapper.Wrap(s.create))
	router.DELETE("/:"+FieldUserID, ginwrapper.Wrap(s.delete))
}

func (s *service) get(c *gin.Context) ginwrapper.Response {
	var req GetPermsRequest
	workspaceID, err := idextractor.ExtractUint64Param(c, consts.FieldWorkspaceId)
	if err != nil {
		return ginwrapper.Response{
			Error: err,
		}
	}
	if err := c.ShouldBind(&req); err != nil {
		return ginwrapper.Response{Error: errors.BadRequest.Wrap(err, "cannot bind request params")}
	}
	response, err := s.repository.GetPermissions(workspaceID, req)
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

func (s *service) create(c *gin.Context) ginwrapper.Response {
	var req CreatePermsRequest
	workspaceID, err := idextractor.ExtractUint64Param(c, consts.FieldWorkspaceId)
	if err != nil {
		return ginwrapper.Response{
			Error: err,
		}
	}
	if err := c.ShouldBind(&req); err != nil {
		return ginwrapper.Response{Error: errors.BadRequest.Wrap(err, "cannot bind request params")}
	}
	err = s.repository.CreatePermissions(workspaceID, req)
	if err != nil {
		return ginwrapper.Response{
			Error: err,
		}
	}
	return ginwrapper.Response{
		Error: errors.Success.NewWithMessage("success"),
	}
}

func (s *service) delete(c *gin.Context) ginwrapper.Response {
	workspaceID, err := idextractor.ExtractUint64Param(c, consts.FieldWorkspaceId)
	if err != nil {
		return ginwrapper.Response{
			Error: err,
		}
	}
	userID, err := idextractor.ExtractUint64Param(c, FieldUserID)
	if err != nil {
		return ginwrapper.Response{
			Error: err,
		}
	}
	logger.Infof("delete user %d at workspace %d", userID, workspaceID)
	err = s.repository.DeletePermission(workspaceID, userID)
	if err != nil {
		return ginwrapper.Response{
			Error: err,
		}
	}
	return ginwrapper.Response{
		Error: errors.Success.NewWithMessage("success"),
	}
}
