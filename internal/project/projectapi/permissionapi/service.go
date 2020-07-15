package permissionapi

import (
	"github.com/gin-gonic/gin"
	"github.com/nkhang/pluto/pkg/errors"
	"github.com/nkhang/pluto/pkg/ginwrapper"
)

type service struct {
	repository Repository
}

const (
	FieldProjectID = "projectId"
)

func NewService(r Repository) *service {
	return &service{
		repository: r,
	}
}

func (s *service) Register(router gin.IRouter) {
	router.POST("", ginwrapper.Wrap(s.createPerm))
	router.GET("", ginwrapper.Wrap(s.getPermissions))
	router.PUT("", ginwrapper.Wrap(s.updatePermission))
}

func (s *service) getPermissions(c *gin.Context) ginwrapper.Response {
	id := c.GetInt64(FieldProjectID)
	resp, err := s.repository.GetList(uint64(id))
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

func (s *service) createPerm(c *gin.Context) ginwrapper.Response {
	id := c.GetInt64(FieldProjectID)
	var req CreatePermRequest
	if err := c.ShouldBind(&req); err != nil {
		return ginwrapper.Response{
			Error: errors.BadRequest.NewWithMessage("error binding request params"),
		}
	}
	err := s.repository.Create(uint64(id), req)
	if err != nil {
		return ginwrapper.Response{
			Error: err,
		}
	}
	return ginwrapper.Response{
		Error: errors.Success.NewWithMessage("success"),
	}
}

func (s *service) updatePermission(c *gin.Context) ginwrapper.Response {
	id := c.GetInt64(FieldProjectID)
	var req UpdatePermissionRequest
	if err := c.ShouldBind(&req); err != nil {
		return ginwrapper.Response{
			Error: errors.BadRequest.NewWithMessage("cannot bind params"),
		}
	}
	perm, err := s.repository.Update(uint64(id), req)
	if err != nil {
		return ginwrapper.Response{
			Error: err,
		}
	}
	return ginwrapper.Response{
		Error: errors.Success.NewWithMessage("success"),
		Data:  perm,
	}
}
