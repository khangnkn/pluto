package pgin

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"

	"github.com/dgrijalva/jwt-go"

	"github.com/nkhang/pluto/pkg/errors"
	"github.com/nkhang/pluto/pkg/ginwrapper"
	"github.com/nkhang/pluto/pkg/logger"

	"github.com/gin-gonic/gin"
)

type payload struct {
	jwt.StandardClaims
	UserID   uint64 `json:"id"`
	Username string `json:"username"`
}

func ApplyVerifyToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqToken := c.GetHeader("Authorization")
		logger.Infof("Authorization header: %s", reqToken)
		splitToken := strings.Split(reqToken, "Bearer ")
		if len(splitToken) < 2 {
			report(c, "authorization header not found")
			return
		}
		key := viper.GetString("jwt.secret")
		logger.Info(key)
		tokenString := splitToken[1]
		tok := payload{}
		token, err := jwt.ParseWithClaims(tokenString, &tok, func(token *jwt.Token) (interface{}, error) {
			return []byte(key), nil
		})
		if err != nil {
			logger.Error(err)
			report(c, "cannot verify user from token")
			return
		}
		if claims, ok := token.Claims.(*payload); ok && token.Valid {
			fmt.Printf("%v %v", claims.UserID, claims.StandardClaims.ExpiresAt)
		} else {
			fmt.Println(err)
		}
		c.Next()
	}
}

func report(c *gin.Context, msg string) {
	ginwrapper.Wrap(func(c *gin.Context) ginwrapper.Response {
		return ginwrapper.Response{
			Error: errors.Unauthorize.NewWithMessage(msg),
		}
	})(c)
}
