package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"go-template/pkg/common/apierrors"
	"go-template/pkg/common/env"
	"go-template/pkg/common/utils"
	"go-template/pkg/handlers/user"
	"strings"
)

func EnforceAuthentication(redisClient *redis.Client, userRepository user.Repository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := strings.ReplaceAll(ctx.GetHeader("Authorization"), "Bearer ", "")
		if tokenString == "" {
			ctx.Error(apierrors.ErrorUnauthorized)
			ctx.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(env.JWT_SECRET), nil
		})
		if errors.Is(err, jwt.ErrTokenMalformed) {
			ctx.Error(apierrors.ErrorUnauthorized)
			ctx.Abort()
			return
		}
		if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
			ctx.Error(apierrors.ErrorUnauthorized)
			ctx.Abort()
			return
		}

		userId := token.Claims.(jwt.MapClaims)["userId"]
		_, err = utils.GetUserById(uint(userId.(float64)), redisClient, userRepository)
		if err != nil {
			ctx.Error(apierrors.ErrorUnauthorized)
			ctx.Abort()
			return
		}

		ctx.Set("userId", token.Claims.(jwt.MapClaims)["userId"])
		ctx.Next()
	}
}
