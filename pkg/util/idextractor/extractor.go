package idextractor

import (
	"github.com/gin-gonic/gin"
	"github.com/nkhang/pluto/pkg/errors"
	"github.com/spf13/cast"
)

func ExtractUint64Param(c *gin.Context, key string) (uint64, error) {
	val := c.Param(key)
	res, err := cast.ToUint64E(val)
	if err != nil {
		return 0, errors.BadRequest.NewWithMessageF("error getting %s from path", val)
	}
	return res, nil
}
