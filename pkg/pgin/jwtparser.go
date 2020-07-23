package pgin

import (
	"strings"

	"github.com/nkhang/pluto/pkg/logger"

	"github.com/spf13/viper"

	"github.com/dgrijalva/jwt-go"

	"github.com/gin-gonic/gin"
	"github.com/nkhang/pluto/pkg/errors"
	"github.com/nkhang/pluto/pkg/ginwrapper"
)

type payload struct {
	jwt.StandardClaims
	UserID   int64  `json:"id"`
	Username string `json:"username"`
}

const (
	FieldUserID = "userId"
)

func ApplyVerifyToken() gin.HandlerFunc {
	key := viper.GetString("jwt.secret")
	logger.Infof("apply verify token middleware with key %s", key)
	return func(c *gin.Context) {
		reqToken := c.GetHeader("Authorization")
		splitToken := strings.Split(reqToken, "Bearer ")
		if len(splitToken) < 2 {
			report(c, "authorization header not found")
			return
		}
		key := viper.GetString("jwt.secret")
		tokenString := splitToken[1]
		logger.Infof("verifying token %s with secret %s", tokenString, key)
		tok := payload{}
		token, err := jwt.ParseWithClaims(tokenString, &tok, func(token *jwt.Token) (interface{}, error) {
			return []byte(key), nil
		})
		if err != nil {
			logger.Error("error verify credential: ", err.Error())
			report(c, "cannot verify user credential")
			return
		}
		if claims, ok := token.Claims.(*payload); ok && token.Valid {
			if claims.UserID == 0 {
				report(c, "error user_id is not valid")
				return
			}
			c.Set(FieldUserID, claims.UserID)
			logger.Infof("credential passed for user %d", claims.UserID)
		} else {
			report(c, "claim payload is invalid")
			return
		}
		c.Next()
	}
}

func ExtractUserIDFromContext(c *gin.Context) uint64 {
	userID := uint64(c.GetInt64(FieldUserID))
	return userID
}

func report(c *gin.Context, msg string) {
	ginwrapper.Wrap(func(c *gin.Context) ginwrapper.Response {
		return ginwrapper.Response{
			Error: errors.Unauthorized.NewWithMessage(msg),
		}
	})(c)
}
