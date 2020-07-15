package projectapi

import (
	"net/http"

	"github.com/nkhang/pluto/internal/workspace/workspaceapi/consts"

	"github.com/nkhang/pluto/pkg/util/paging"

	"github.com/nkhang/pluto/internal/project/projectapi/permissionapi"

	"github.com/nkhang/pluto/pkg/pgin"

	"github.com/nkhang/pluto/internal/project"

	"github.com/gin-gonic/gin"
	"github.com/nkhang/pluto/pkg/util/idextractor"

	"github.com/nkhang/pluto/pkg/errors"
	"github.com/nkhang/pluto/pkg/ginwrapper"
)

type service struct {
	repository       Repository
	projectRepo      project.Repository
	permissionRouter pgin.Router
}

const (
	FieldProjectID = "projectId"
)

func NewService(r Repository, projectRepo project.Repository) *service {
	permRepo := permissionapi.NewProjectPermissionAPIRepository(projectRepo)
	permRouter := permissionapi.NewService(permRepo)
	return &service{
		repository:       r,
		projectRepo:      projectRepo,
		permissionRouter: permRouter,
	}
}

func (s *service) Register(router gin.IRouter) {
	router.POST("", ginwrapper.Wrap(s.create))
	router.GET("", ginwrapper.Wrap(s.getForWorkspace))
	detailRouter := router.Group("/:"+FieldProjectID, s.verifyProjectIDMdw())
	{
		detailRouter.GET("", ginwrapper.Wrap(s.get))
		detailRouter.PUT("", ginwrapper.Wrap(s.update))
		detailRouter.DELETE("", ginwrapper.Wrap(s.delete))
	}
	s.permissionRouter.Register(detailRouter.Group("/perms"))
}

func (s *service) RegisterStandalone(router gin.IRouter) {
	router.GET("", ginwrapper.Wrap(s.getForUser))
}

func (s *service) getForUser(c *gin.Context) ginwrapper.Response {
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

func (s *service) getForWorkspace(c *gin.Context) ginwrapper.Response {
	id := uint64(c.GetInt64(consts.FieldWorkspaceId))
	var pg paging.Paging
	if err := c.ShouldBindQuery(&pg); err != nil {
		return ginwrapper.Response{
			Error: errors.BadRequest.NewWithMessage("cannot bind object"),
		}
	}
	resp, err := s.repository.GetForWorkspace(id, pg)
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

func (s *service) create(c *gin.Context) ginwrapper.Response {
	workspaceID := uint64(c.GetInt64(consts.FieldWorkspaceId))
	var req CreateProjectRequest
	if err := c.ShouldBind(&req); err != nil {
		return ginwrapper.Response{
			Error: errors.BadRequest.Wrap(err, "cannot bind request params"),
		}
	}
	resp, err := s.repository.Create(workspaceID, req)
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
