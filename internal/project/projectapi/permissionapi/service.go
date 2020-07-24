package permissionapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nkhang/pluto/internal/project"
	"github.com/nkhang/pluto/internal/project/projectapi"
	"github.com/nkhang/pluto/pkg/errors"
	"github.com/nkhang/pluto/pkg/ginwrapper"
	"github.com/nkhang/pluto/pkg/util/idextractor"
)

type service struct {
	repository  Repository
	projectRepo project.Repository
}

const (
	FieldUserID = "userId"
)

func NewService(r Repository, pr project.Repository) *service {
	return &service{
		repository:  r,
		projectRepo: pr,
	}
}

func (s *service) Register(router gin.IRouter) {
	router.POST("", ginwrapper.Wrap(s.createPerm))
	router.GET("", ginwrapper.Wrap(s.getPermissions))
	router.PUT("", ginwrapper.Wrap(s.updatePermission))
	router.DELETE("/:"+FieldUserID, s.verifyPermission(), ginwrapper.Wrap(s.delete))
}

func (s *service) getPermissions(c *gin.Context) ginwrapper.Response {
	id := c.GetInt64(projectapi.FieldProjectID)
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
	id := c.GetInt64(projectapi.FieldProjectID)
	var req CreatePermRequest
	if err := c.ShouldBind(&req); err != nil {
		return ginwrapper.Response{
			Error: errors.BadRequest.NewWithMessage("error binding request params"),
		}
	}
	prj, err := s.repository.Create(uint64(id), req)
	if err != nil {
		return ginwrapper.Response{
			Error: err,
		}
	}
	return ginwrapper.Response{
		Error: errors.Success.NewWithMessage("success"),
		Data:  prj,
	}
}

func (s *service) updatePermission(c *gin.Context) ginwrapper.Response {
	id := c.GetInt64(projectapi.FieldProjectID)
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

func (s *service) delete(c *gin.Context) ginwrapper.Response {
	projectID := uint64(c.GetInt64(projectapi.FieldProjectID))
	userID, err := idextractor.ExtractUint64Param(c, FieldUserID)
	if err != nil {
		return ginwrapper.Response{
			Error: err,
		}
	}
	err = s.repository.Delete(projectID, userID)
	if err != nil {
		return ginwrapper.Response{
			Error: err,
		}
	}
	return ginwrapper.Response{
		Error: errors.Success.NewWithMessage("success"),
	}
}

func (s *service) verifyPermission() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		projectID := uint64(c.GetInt64(projectapi.FieldProjectID))
		userID, err := idextractor.ExtractUint64Param(c, FieldUserID)
		if err != nil {
			err := errors.BadRequest.NewWithMessageF("user id %d not found", userID)
			ginwrapper.Report(c, http.StatusOK, err, nil)
			return
		}
		if userID == 0 {
			err := errors.BadRequest.NewWithMessage("user id ID must be other than 0")
			ginwrapper.Report(c, http.StatusOK, err, nil)
			return
		}
		p, err := s.projectRepo.GetPermission(userID, projectID)
		if err != nil {
			ginwrapper.Report(c, http.StatusOK, err, nil)
			return
		}
		if p.ProjectID != projectID {
			err = errors.ProjectNotFound.NewWithMessageF("user %d does not belong to project %d", userID, projectID)
			ginwrapper.Report(c, http.StatusOK, err, nil)
			return
		}
		c.Set(FieldUserID, userID)
		c.Next()
	})
}
