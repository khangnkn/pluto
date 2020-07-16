package datasetapi

import (
	"net/http"

	"github.com/nkhang/pluto/pkg/pgin"

	"github.com/gin-gonic/gin"
	"github.com/nkhang/pluto/internal/dataset"
	"github.com/nkhang/pluto/internal/project/projectapi"
	"github.com/nkhang/pluto/pkg/util/idextractor"

	"github.com/nkhang/pluto/pkg/errors"
	"github.com/nkhang/pluto/pkg/ginwrapper"
	"github.com/nkhang/pluto/pkg/logger"
)

const (
	FieldDatasetID = "datasetId"
)

type service struct {
	repository  Repository
	datasetRepo dataset.Repository
	imageRouter pgin.Router
}

func NewService(r Repository, datasetRepo dataset.Repository, imageRouter pgin.Router) *service {
	return &service{
		repository:  r,
		datasetRepo: datasetRepo,
		imageRouter: imageRouter,
	}
}

func (s *service) Register(router gin.IRouter) {
	router.GET("", ginwrapper.Wrap(s.getByProjectID))
	router.POST("", ginwrapper.Wrap(s.create))
	detailRouter := router.Group("/:"+FieldDatasetID, s.verifyDataset())
	{
		detailRouter.POST("")
		detailRouter.GET("", ginwrapper.Wrap(s.getByID))
		detailRouter.DELETE("", ginwrapper.Wrap(s.del))
		detailRouter.GET("/link", ginwrapper.Wrap(s.getLink))
	}
	s.imageRouter.Register(detailRouter.Group("/images"))
}

func (s *service) getByID(c *gin.Context) ginwrapper.Response {
	datasetID := uint64(c.GetInt64(FieldDatasetID))
	dataset, err := s.repository.GetByID(datasetID)
	if err != nil {
		return ginwrapper.Response{
			Error: err,
		}
	}
	return ginwrapper.Response{
		Error: errors.Success.NewWithMessage("Success"),
		Data:  dataset,
	}
}

func (s *service) getByProjectID(c *gin.Context) ginwrapper.Response {
	projectID := uint64(c.GetInt64(projectapi.FieldProjectID))
	datasets, err := s.repository.GetByProjectID(projectID)
	if err != nil {
		return ginwrapper.Response{
			Error: err,
		}
	}
	return ginwrapper.Response{
		Error: errors.Success.NewWithMessage("Success"),
		Data:  datasets,
	}
}

func (s *service) create(c *gin.Context) ginwrapper.Response {
	projectID := uint64(c.GetInt64(projectapi.FieldProjectID))
	var req CreateDatasetRequest
	if err := c.ShouldBind(&req); err != nil {
		return ginwrapper.Response{
			Error: errors.BadRequest.Wrap(err, "cannot bind request"),
		}
	}
	dataset, err := s.repository.CreateDataset(req.Title, req.Description, projectID)
	if err != nil {
		return ginwrapper.Response{
			Error: err,
		}
	}
	return ginwrapper.Response{
		Error: errors.Success.NewWithMessage("success"),
		Data:  dataset,
	}
}

func (s *service) clone(c *gin.Context) ginwrapper.Response {
	projectID := uint64(c.GetInt64(projectapi.FieldProjectID))
	logger.Info("start cloning project")
	var req CloneDatasetRequest
	if err := c.ShouldBind(&req); err != nil {
		return ginwrapper.Response{
			Error: errors.BadRequest.NewWithMessage("cannot bind clone dataset request"),
		}
	}
	cloned, err := s.repository.CloneDataset(projectID, req.Token)
	if err != nil {
		return ginwrapper.Response{
			Error: err,
		}
	}
	return ginwrapper.Response{
		Error: errors.Success.NewWithMessage("success"),
		Data:  cloned,
	}
}

func (s *service) del(c *gin.Context) ginwrapper.Response {
	datasetID := uint64(c.GetInt64(FieldDatasetID))
	err := s.repository.Delete(datasetID)
	if err != nil {
		return ginwrapper.Response{Error: err}
	}
	return ginwrapper.Response{
		Error: errors.Success.NewWithMessage("success"),
	}
}

func (s *service) verifyDataset() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		param := c.Param(FieldDatasetID)
		if param == "parse" {
			ginwrapper.Wrap(s.parseLink)(c)
			return
		}
		if param == "clone" {
			ginwrapper.Wrap(s.clone)(c)
			return
		}
		datasetID, err := idextractor.ExtractInt64Param(c, FieldDatasetID)
		if err != nil {
			err := errors.BadRequest.NewWithMessageF("dataset %d not found", datasetID)
			ginwrapper.Report(c, http.StatusOK, err, nil)
			return
		}
		if datasetID == 0 {
			err := errors.BadRequest.NewWithMessage("dataset ID must be other than 0")
			ginwrapper.Report(c, http.StatusOK, err, nil)
			return
		}
		p, err := s.datasetRepo.Get(uint64(datasetID))
		if err != nil {
			ginwrapper.Report(c, http.StatusOK, err, nil)
			return
		}
		if projectID := uint64(c.GetInt64(projectapi.FieldProjectID)); p.ProjectID != projectID {
			err = errors.ProjectNotFound.NewWithMessageF("dataset %d does not belong to project %d", datasetID, projectID)
			ginwrapper.Report(c, http.StatusOK, err, nil)
			return
		}
		c.Set(FieldDatasetID, datasetID)
		c.Next()
	})
}

func (s *service) getLink(c *gin.Context) ginwrapper.Response {
	datasetID := uint64(c.GetInt64(FieldDatasetID))
	url, err := s.repository.GetLink(datasetID)
	if err != nil {
		return ginwrapper.Response{
			Error: err,
		}
	}
	return ginwrapper.Response{
		Error: errors.Success.NewWithMessage("success"),
		Data:  url,
	}
}

func (s *service) parseLink(c *gin.Context) ginwrapper.Response {
	var req ParseLinkRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		return ginwrapper.Response{
			Error: errors.BadRequest.NewWithMessage("error getting link to parse"),
		}
	}
	resp, err := s.repository.ParseLink(req.Link)
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
