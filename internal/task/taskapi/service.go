package taskapi

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nkhang/pluto/internal/task"
	"github.com/nkhang/pluto/pkg/errors"
	"github.com/nkhang/pluto/pkg/ginwrapper"
	"github.com/nkhang/pluto/pkg/logger"
	"github.com/nkhang/pluto/pkg/util/idextractor"
)

const (
	fieldTaskID       = "taskId"
	fieldTaskDetailID = "taskDetailId"
)

type service struct {
	repository Repository
	taskRepo   task.Repository
}

func NewService(r Repository, tr task.Repository) *service {
	return &service{
		repository: r,
		taskRepo:   tr,
	}
}

func (s *service) RegisterStandalone(router gin.IRouter) {
	router.GET("", ginwrapper.Wrap(s.getList))
	router.POST("", ginwrapper.Wrap(s.createTask))
	detailRouter := router.Group("/:"+fieldTaskID, s.verifyTask())
	{
		detailRouter.DELETE("", ginwrapper.Wrap(s.delete))
		detailRouter.GET("", ginwrapper.Wrap(s.get))
		detailRouter.PUT("/details/:"+fieldTaskDetailID, ginwrapper.Wrap(s.updateTaskDetail))
		detailRouter.GET("/details", ginwrapper.Wrap(s.getTaskDetails))
	}
}

func (s *service) createTask(c *gin.Context) ginwrapper.Response {
	var req CreateTaskRequest
	if err := c.ShouldBind(&req); err != nil {
		return ginwrapper.Response{
			Error: errors.BadRequest.NewWithMessage("error binding create task request"),
		}
	}
	err := s.repository.CreateTask(req)
	if err != nil {
		return ginwrapper.Response{
			Error: err,
		}
	}
	return ginwrapper.Response{
		Error: errors.Success.NewWithMessage("success"),
	}
}

func (s *service) get(c *gin.Context) ginwrapper.Response {
	taskID := uint64(c.GetInt64(fieldTaskID))
	resp, err := s.repository.GetTask(taskID)
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

func (s *service) getList(c *gin.Context) ginwrapper.Response {
	var request GetTasksRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		logger.Error(err)
		return ginwrapper.Response{
			Error: errors.BadRequest.NewWithMessage("error binding get tasks request"),
		}
	}
	response, err := s.repository.GetTasks(request)
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

func (s *service) delete(c *gin.Context) ginwrapper.Response {
	taskID := uint64(c.GetInt64(fieldTaskID))
	err := s.repository.DeleteTask(uint64(taskID))
	if err != nil {
		return ginwrapper.Response{
			Error: err,
		}
	}
	return ginwrapper.Response{
		Error: errors.Success.NewWithMessage("success"),
	}
}

func (s *service) getTaskDetails(c *gin.Context) ginwrapper.Response {
	taskID := uint64(c.GetInt64(fieldTaskID))
	var req GetTaskDetailsRequest
	if err := c.ShouldBind(&req); err != nil {
		return ginwrapper.Response{
			Error: errors.BadRequest.NewWithMessage("cannot bind request params"),
		}
	}
	details, err := s.repository.GetTaskDetails(taskID, req)
	if err != nil {
		return ginwrapper.Response{Error: err}
	}
	return ginwrapper.Response{
		Error: errors.Success.NewWithMessage("success"),
		Data:  details,
	}
}

func (s *service) updateTaskDetail(c *gin.Context) ginwrapper.Response {
	b, err := ioutil.ReadAll(c.Request.Body)
	if err != nil || len(b) == 0 {
		return ginwrapper.Response{
			Error: errors.BadRequest.NewWithMessage("cannot read request body"),
		}
	}
	var request UpdateTaskDetailRequest
	detailID, err := idextractor.ExtractUint64Param(c, fieldTaskDetailID)
	if err != nil {
		return ginwrapper.Response{Error: err}
	}
	taskID := uint64(c.GetInt64(fieldTaskID))
	if err := json.Unmarshal(b, &request); err != nil {
		logger.Error(err)
		return ginwrapper.Response{
			Error: errors.BadRequest.Wrap(err, "error binding update task detail request"),
		}
	}
	logger.Infof("request %+v", request)
	response, err := s.repository.UpdateTaskDetail(taskID, detailID, request)
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

func (s *service) verifyTask() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		taskID, err := idextractor.ExtractInt64Param(c, fieldTaskID)
		if err != nil {
			err := errors.BadRequest.NewWithMessageF("task %d not found", taskID)
			ginwrapper.Report(c, http.StatusOK, err, nil)
			return
		}
		if taskID == 0 {
			err := errors.BadRequest.NewWithMessage("task ID must be other than 0")
			ginwrapper.Report(c, http.StatusOK, err, nil)
			return
		}
		_, err = s.taskRepo.GetTask(uint64(taskID))
		if err != nil {
			ginwrapper.Report(c, http.StatusOK, err, nil)
			return
		}
		c.Set(fieldTaskID, taskID)
		c.Next()
	})
}
