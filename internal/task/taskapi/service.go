package taskapi

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/nkhang/pluto/pkg/pgin"

	"github.com/nats-io/nats.go"
	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"
	"github.com/nkhang/pluto/internal/project/projectapi"
	"github.com/nkhang/pluto/internal/task"
	"github.com/nkhang/pluto/pkg/errors"
	"github.com/nkhang/pluto/pkg/ginwrapper"
	"github.com/nkhang/pluto/pkg/logger"
	"github.com/nkhang/pluto/pkg/util/idextractor"
)

const (
	FieldTaskID       = "taskId"
	fieldTaskDetailID = "taskDetailId"
)

type Service struct {
	repository  Repository
	taskRepo    task.Repository
	statsRouter pgin.Router
}

func NewService(r Repository, tr task.Repository, statsRouter pgin.Router) *Service {
	return &Service{
		repository:  r,
		taskRepo:    tr,
		statsRouter: statsRouter,
	}
}

func (s *Service) Register(router gin.IRouter) {
	router.POST("", ginwrapper.Wrap(s.createTask))
	router.GET("", ginwrapper.Wrap(s.getListForProject))
	detailRouter := router.Group("/:"+FieldTaskID, s.verifyTask())
	{
		detailRouter.DELETE("", ginwrapper.Wrap(s.delete))
		detailRouter.GET("", ginwrapper.Wrap(s.get))
		detailRouter.GET("/details", ginwrapper.Wrap(s.getTaskDetails))
	}
	s.statsRouter.Register(detailRouter)
}

func (s *Service) RegisterNATS(ec *nats.EncodedConn) error {
	var topic = viper.GetString("annotation.updatetask")
	logger.Info(topic)
	_, err := ec.Subscribe(topic, s.handleUpdateTask)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) RegisterStandalone(router gin.IRouter) {
	router.GET("", ginwrapper.Wrap(s.getListForUser))
}

func (s *Service) RegisterInternal(router gin.IRouter) {
	taskRouter := router.Group("/workspaces/:workspaceId/projects/:projectId/tasks/:"+FieldTaskID, s.verifyTask())
	{
		taskRouter.PUT("/details/:"+fieldTaskDetailID, ginwrapper.Wrap(s.updateTaskDetail))
	}
}

func (s *Service) handleUpdateTask(msg *nats.Msg) {
	b := msg.Data
	logger.Infof("%s", b)
	var req NATSUpdateDetailRequest
	err := json.Unmarshal(msg.Data, &req)
	if err != nil {
		logger.Errorf("error unmarshal message from nats. error %v. msg %s", err, msg.Data)
		return
	}
	_, err = s.repository.UpdateTaskDetail(req.TaskID, req.DetailID, UpdateTaskDetailRequest{Status: req.Status})
	if err != nil {
		logger.Infof("error updating task detail. task %d, detail %d, status %d, err %v", req.TaskID, req.DetailID, req.Status, err)
		return
	}
	logger.Infof("update task detail successfully. task %d, detail %d, status %d", req.TaskID, req.DetailID, req.Status)
}

func (s *Service) createTask(c *gin.Context) ginwrapper.Response {
	projectID := uint64(c.GetInt64(projectapi.FieldProjectID))
	assigner := pgin.ExtractUserIDFromContext(c)
	var req CreateTaskRequest
	if err := c.ShouldBind(&req); err != nil {
		return ginwrapper.Response{
			Error: errors.BadRequest.NewWithMessage("error binding create task request"),
		}
	}
	err := s.repository.CreateTask(projectID, assigner, req)
	if err != nil {
		return ginwrapper.Response{
			Error: err,
		}
	}
	return ginwrapper.Response{
		Error: errors.Success.NewWithMessage("success"),
	}
}

func (s *Service) get(c *gin.Context) ginwrapper.Response {
	taskID := uint64(c.GetInt64(FieldTaskID))
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

func (s *Service) getListForUser(c *gin.Context) ginwrapper.Response {
	userID := pgin.ExtractUserIDFromContext(c)
	var request GetTasksRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		logger.Error(err)
		return ginwrapper.Response{
			Error: errors.BadRequest.NewWithMessage("error binding get tasks request"),
		}
	}
	response, err := s.repository.GetTasks(userID, request)
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

func (s *Service) getListForProject(c *gin.Context) ginwrapper.Response {
	projectID := uint64(c.GetInt64(projectapi.FieldProjectID))
	userID := pgin.ExtractUserIDFromContext(c)
	var req GetTasksRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Error(err)
		return ginwrapper.Response{
			Error: errors.BadRequest.NewWithMessage("error binding get tasks request"),
		}
	}
	response, err := s.repository.GetTaskForProject(projectID, userID, req)
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

func (s *Service) delete(c *gin.Context) ginwrapper.Response {
	taskID := uint64(c.GetInt64(FieldTaskID))
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

func (s *Service) getTaskDetails(c *gin.Context) ginwrapper.Response {
	taskID := uint64(c.GetInt64(FieldTaskID))
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

func (s *Service) updateTaskDetail(c *gin.Context) ginwrapper.Response {
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
	taskID := uint64(c.GetInt64(FieldTaskID))
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

func (s *Service) verifyTask() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		taskID, err := idextractor.ExtractInt64Param(c, FieldTaskID)
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
		c.Set(FieldTaskID, taskID)
		c.Next()
	})
}
