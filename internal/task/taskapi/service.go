package taskapi

import (
	"github.com/gin-gonic/gin"
	"github.com/nkhang/pluto/pkg/errors"
	"github.com/nkhang/pluto/pkg/ginwrapper"
	"github.com/nkhang/pluto/pkg/logger"
	"github.com/nkhang/pluto/pkg/util/idextractor"
	"github.com/spf13/cast"
)

const (
	fieldTaskID       = "taskId"
	fieldTaskDetailID = "taskDetailId"
)

type service struct {
	repository Repository
}

func NewService(r Repository) *service {
	return &service{repository: r}
}

func (s *service) Register(router gin.IRouter) {
	router.GET("", ginwrapper.Wrap(s.get))
	router.POST("/", ginwrapper.Wrap(s.createTask))
	router.DELETE("/:"+fieldTaskID, ginwrapper.Wrap(s.delete))
	router.PUT("/:"+fieldTaskID+"/details/:"+fieldTaskDetailID, ginwrapper.Wrap(s.updateTaskDetail))
	router.GET("/:"+fieldTaskID+"/details", ginwrapper.Wrap(s.getTaskDetails))
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
	id, err := idextractor.ExtractUint64Param(c, fieldTaskID)
	if err != nil {
		return ginwrapper.Response{
			Error: err,
		}
	}
	err = s.repository.DeleteTask(id)
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
	idStr := c.Param(fieldTaskID)
	taskID, err := cast.ToUint64E(idStr)
	if err != nil {
		return ginwrapper.Response{
			Error: errors.BadRequest.Wrap(err, "invalid task id"),
		}
	}
	var req GetTaskDetailsRequest
	if err := c.ShouldBind(&req); err != nil {
		return ginwrapper.Response{
			Error: errors.BadRequest.NewWithMessage("cannot bind request params"),
		}
	}
	req.TaskID = taskID
	details, err := s.repository.GetTaskDetails(req)
	if err != nil {
		return ginwrapper.Response{Error: err}
	}
	return ginwrapper.Response{
		Error: errors.Success.NewWithMessage("success"),
		Data:  details,
	}
}

func (s *service) updateTaskDetail(c *gin.Context) ginwrapper.Response {
	var request UpdateTaskDetailRequest
	detailID, err := idextractor.ExtractUint64Param(c, fieldTaskDetailID)
	if err != nil {
		return ginwrapper.Response{Error: err}
	}
	taskID, err := idextractor.ExtractUint64Param(c, fieldTaskID)
	if err != nil {
		return ginwrapper.Response{Error: err}
	}
	if err := c.ShouldBind(&request); err != nil {
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
