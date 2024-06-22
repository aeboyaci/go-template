package middlewares

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"go-template/pkg/common/apierrors"
	"go-template/pkg/common/env"
	"go-template/pkg/handlers/user"
	"go-template/pkg/models"
	"gorm.io/gorm"
	"strings"
)

func EnforceAuthentication(databaseClient *gorm.DB, redisClient *redis.Client, userRepository user.Repository) gin.HandlerFunc {
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
		_, err = getUserById(uint(userId.(float64)), databaseClient, redisClient, userRepository)
		if err != nil {
			ctx.Error(apierrors.ErrorUnauthorized)
			ctx.Abort()
			return
		}

		ctx.Set("userId", token.Claims.(jwt.MapClaims)["userId"])
		ctx.Next()
	}
}

func getUserById(userId uint, databaseClient *gorm.DB, redisClient *redis.Client, repository user.Repository) (*models.User, error) {
	ctx := context.Background()

	userKey := fmt.Sprintf("user:%d", userId)
	userJson, err := redisClient.Get(ctx, userKey).Result()
	if errors.Is(err, redis.Nil) {
		user, err := repository.FindUserById(databaseClient, userId)
		if err != nil {
			return nil, err
		}

		userJson, err := json.Marshal(user)
		if err != nil {
			return nil, err
		}
		err = redisClient.Set(ctx, userKey, userJson, 0).Err()
		if err != nil {
			return nil, err
		}

		return &user, nil
	} else if err != nil {
		return nil, err
	}

	var user models.User
	err = json.Unmarshal([]byte(userJson), &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
