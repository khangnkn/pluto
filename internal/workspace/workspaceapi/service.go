package workspaceapi

import (
	"net/http"

	"github.com/nkhang/pluto/internal/workspace/workspaceapi/permissionapi"

	"github.com/gin-gonic/gin"
	"github.com/nkhang/pluto/internal/workspace"
	pgin "github.com/nkhang/pluto/pkg/pgin"
	"github.com/nkhang/pluto/pkg/util/idextractor"
	"github.com/spf13/cast"

	"github.com/nkhang/pluto/pkg/errors"
	"github.com/nkhang/pluto/pkg/ginwrapper"
)

const (
	FieldWorkspaceID = "workspaceId"
)

type service struct {
	repository    Repository
	workspaceRepo workspace.Repository
	permRouter    pgin.Router
	projectRouter pgin.Router
}

func NewService(r Repository,
	workspaceRepo workspace.Repository, pr pgin.Router) *service {
	permRepo := permissionapi.NewRepository(workspaceRepo)
	permRouter := permissionapi.NewService(permRepo)
	return &service{
		repository:    r,
		workspaceRepo: workspaceRepo,
		permRouter:    permRouter,
		projectRouter: pr,
	}
}

func (s *service) RegisterStandalone(router gin.IRouter) {
	router.GET("", ginwrapper.Wrap(s.getByUserID))
	router.POST("", ginwrapper.Wrap(s.create))
	detailRouter := router.Group("/:"+FieldWorkspaceID, s.verifyWorkspace())
	{
		detailRouter.GET("", ginwrapper.Wrap(s.get))
		detailRouter.PUT("", ginwrapper.Wrap(s.update))
		detailRouter.DELETE("", ginwrapper.Wrap(s.delete))
	}
	s.permRouter.Register(detailRouter.Group("/perms"))
	s.projectRouter.Register(detailRouter.Group("/projects"))
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

func (s *service) update(c *gin.Context) ginwrapper.Response {
	var req UpdateWorkspaceRequest
	workspaceID, err := idextractor.ExtractUint64Param(c, FieldWorkspaceID)
	if err != nil {
		return ginwrapper.Response{
			Error: err,
		}
	}
	if err := c.ShouldBind(&req); err != nil {
		return ginwrapper.Response{
			Error: errors.BadRequest.Wrap(err, "cannot bind update request"),
		}
	}
	w, err := s.repository.UpdateWorkspace(workspaceID, req)
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
	workspaceID, err := idextractor.ExtractUint64Param(c, FieldWorkspaceID)
	if err != nil {
		return ginwrapper.Response{
			Error: err,
		}
	}
	err = s.repository.DeleteWorkspace(workspaceID)
	if err != nil {
		return ginwrapper.Response{
			Error: err,
		}
	}
	return ginwrapper.Response{
		Error: errors.Success.NewWithMessage("success"),
	}
}

func (s *service) verifyWorkspace() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		workspaceID, err := idextractor.ExtractInt64Param(c, FieldWorkspaceID)
		if err != nil {
			err := errors.BadRequest.NewWithMessageF("workspace ID %d is invalid", workspaceID)
			ginwrapper.Report(c, http.StatusOK, err, nil)
			return
		}
		if workspaceID == 0 {
			err := errors.BadRequest.NewWithMessage("workspace ID must be other than 0")
			ginwrapper.Report(c, http.StatusOK, err, nil)
			return
		}
		_, err = s.workspaceRepo.Get(uint64(workspaceID))
		if err != nil {
			ginwrapper.Report(c, http.StatusOK, err, nil)
			return
		}
		c.Set(FieldWorkspaceID, workspaceID)
		c.Next()
	})
}
