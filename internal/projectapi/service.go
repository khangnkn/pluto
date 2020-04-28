package projectapi

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/nkhang/pluto/pkg/errors"
	"github.com/nkhang/pluto/pkg/ginwrapper"
)

type service struct {
	repository Repository
}

const (
	fieldID = "id"
)

func NewService(r Repository) *service {
	return &service{repository: r}
}

func (s *service) Register(router gin.IRouter) {
	router.GET("/:"+fieldID, ginwrapper.Wrap(s.get))
}

func (s *service) get(c *gin.Context) ginwrapper.Response {
	idStr := c.Param(fieldID)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return ginwrapper.Response{
			Error: errors.BadRequest.Wrap(err, "cannot get id"),
		}
	}
	if id <= 0 {
		return ginwrapper.Response{
			Error: errors.BadRequest.NewWithMessage("id must greater than 0"),
		}
	}
	p, err := s.repository.GetByID(uint64(id))
	if err != nil {
		return ginwrapper.Response{
			Error: err,
		}
	}
	return ginwrapper.Response{
		Error: errors.Success.NewWithMessage("success"),
		Data:  p,
	}
}
