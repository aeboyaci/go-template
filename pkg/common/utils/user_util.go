package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"go-template/pkg/handlers/user"
	"go-template/pkg/models"
)

func GetUserById(userId uint, redisClient *redis.Client, repository user.Repository) (*models.User, error) {
	ctx := context.Background()

	userKey := fmt.Sprintf("user:%d", userId)
	userJson, err := redisClient.Get(ctx, userKey).Result()
	if errors.Is(err, redis.Nil) {
		user, err := repository.FindUserById(userId)
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
