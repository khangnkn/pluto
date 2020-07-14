package ginwrapper

import (
	"github.com/gin-gonic/gin"
	"github.com/nkhang/pluto/pkg/errors"
)

type Response struct {
	HttpCode int
	Error    error
	Data     interface{}
}

type response struct {
	ReturnCode    int         `json:"status"`
	ReturnMessage string      `json:"msg"`
	Data          interface{} `json:"data,omitempty"`
}

type GinHandlerFunc func(c *gin.Context) Response

func Wrap(fn GinHandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		r := fn(c)
		var httpCode int = 200
		if r.HttpCode != 0 {
			httpCode = r.HttpCode
		}
		Report(c, httpCode, r.Error, r.Data)
	}
}

func Report(c *gin.Context, code int, err error, data interface{}) {
	e, ok := err.(errors.CustomError)
	if !ok {
		e = errors.Unknown.NewWithMessage("Unknown error")
	}
	returnObj := response{
		ReturnCode:    int(e.Code),
		ReturnMessage: e.Message,
		Data:          data,
	}
	c.AbortWithStatusJSON(code, returnObj)
}
