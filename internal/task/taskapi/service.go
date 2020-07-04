package taskapi

import (
	"github.com/gin-gonic/gin"
	"github.com/nkhang/pluto/pkg/errors"
	"github.com/nkhang/pluto/pkg/ginwrapper"
	"github.com/spf13/cast"
)

const fieldTaskID = "task_id"

type service struct {
	repository Repository
}

func NewService(r Repository) *service {
	return &service{repository: r}
}

func (s *service) Register(router gin.IRouter) {
	router.POST("/", ginwrapper.Wrap(s.createTask))
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
