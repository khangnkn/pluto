package pgin

import (
	"strings"

	"github.com/nkhang/pluto/pkg/errors"
	"github.com/nkhang/pluto/pkg/ginwrapper"
	"github.com/nkhang/pluto/pkg/logger"

	"github.com/gin-gonic/gin"
)

type payload struct {
	UserID   uint64 `json:"id"`
	Username string `json:"username"`
}

func ApplyVerifyToken() gin.HandlerFunc {
	return func(c *gin.Context) {

		logger.Info(c.GetHeader("Authorization"))
		logger.Info(c.GetHeader("Bearer"))
		logger.Info(c.Request.Header)
		reqToken := c.GetHeader("Authorization")
		splitToken := strings.Split(reqToken, "Bearer ")
		if len(splitToken) < 2 {
			ginwrapper.Wrap(func(c *gin.Context) ginwrapper.Response {
				return ginwrapper.Response{
					Error: errors.Unauthorize.NewWithMessage("Unauthorized"),
				}
			})(c)
			return
		}
		reqToken = splitToken[1]
		logger.Info(reqToken)
		c.Next()
	}
}
