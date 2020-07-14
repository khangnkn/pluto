package projectapi

import (
	"net/http"

	"github.com/nkhang/pluto/internal/project"

	"github.com/gin-gonic/gin"
	"github.com/nkhang/pluto/pkg/util/idextractor"

	"github.com/nkhang/pluto/pkg/errors"
	"github.com/nkhang/pluto/pkg/ginwrapper"
)

type service struct {
	repository  Repository
	projectRepo project.Repository
}

const (
	FieldProjectID = "projectId"
)

func NewService(r Repository, projectRepo project.Repository) *service {
	return &service{
		repository:  r,
		projectRepo: projectRepo,
	}
}

func (s *service) Register(router gin.IRouter) {
	router.GET("", ginwrapper.Wrap(s.getAll))
	router.POST("", ginwrapper.Wrap(s.create))
	detailRouter := router.Group("/:" + FieldProjectID).Use(s.verifyProjectIDMdw())
	{
		detailRouter.GET("", ginwrapper.Wrap(s.get))
		detailRouter.PUT("", ginwrapper.Wrap(s.update))
		detailRouter.DELETE("", ginwrapper.Wrap(s.delete))
		detailRouter.POST("/perm", ginwrapper.Wrap(s.createPerm))
		detailRouter.GET("/perm", ginwrapper.Wrap(s.getPermissions))
		detailRouter.PUT("/perm", ginwrapper.Wrap(s.updatePermission))
	}
}

func (s *service) getAll(c *gin.Context) ginwrapper.Response {
	var req GetProjectRequest
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
	id := c.GetInt64(FieldProjectID)
	p, err := s.repository.GetByID(uint64(id))
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

func (s *service) getPermissions(c *gin.Context) ginwrapper.Response {
	id := c.GetInt64(FieldProjectID)
	resp, err := s.repository.GetPermissions(uint64(id))
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

func (s *service) create(c *gin.Context) ginwrapper.Response {
	var req CreateProjectRequest
	if err := c.ShouldBind(&req); err != nil {
		return ginwrapper.Response{
			Error: errors.BadRequest.Wrap(err, "cannot bind request params"),
		}
	}
	resp, err := s.repository.Create(req)
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
	var req CreatePermParams
	if err := c.ShouldBind(&req); err != nil {
		return ginwrapper.Response{
			Error: errors.BadRequest.NewWithMessage("error binding request params"),
		}
	}
	req.ProjectID = uint64(id)
	err := s.repository.CreatePerm(req)
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
	perm, err := s.repository.UpdatePerm(uint64(id), req)
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

func (s *service) update(c *gin.Context) ginwrapper.Response {
	var req UpdateProjectRequest
	id := c.GetInt64(FieldProjectID)
	if err := c.ShouldBind(&req); err != nil {
		return ginwrapper.Response{
			Error: errors.BadRequest.Wrap(err, "cannot bind update request"),
		}
	}
	w, err := s.repository.UpdateProject(uint64(id), req)
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

func (s *service) delete(c *gin.Context) ginwrapper.Response {
	id := c.GetInt64(FieldProjectID)
	err := s.repository.DeleteProject(uint64(id))
	if err != nil {
		return ginwrapper.Response{
			Error: err,
		}
	}
	return ginwrapper.Response{
		Error: errors.Success.NewWithMessage("success"),
	}
}

func (s *service) verifyProjectIDMdw() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		projectID, err := idextractor.ExtractInt64Param(c, FieldProjectID)
		if err != nil {
			err := errors.BadRequest.NewWithMessageF("project %d not found", projectID)
			ginwrapper.Report(c, http.StatusOK, err, nil)
			return
		}
		if projectID == 0 {
			err := errors.BadRequest.NewWithMessage("project ID must be other than 0")
			ginwrapper.Report(c, http.StatusOK, err, nil)
			return
		}
		_, err = s.projectRepo.Get(uint64(projectID))
		if err != nil {
			ginwrapper.Report(c, http.StatusOK, err, nil)
			return
		}
		c.Set(FieldProjectID, projectID)
		c.Next()
	})
}
