package redis

import (
	"context"
	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"go-template/pkg/common/env"
	"sync"
)

var (
	client *redis.Client
	once   sync.Once
)

func Initialize() error {
	if env.MODE != "production" {
		s, err := miniredis.Run()
		if err != nil {
			return err
		}

		env.REDIS_URL = s.Addr()
	}

	client = redis.NewClient(&redis.Options{
		Addr:     env.REDIS_URL,
		Password: "",
		DB:       0,
	})

	err := client.Ping(context.Background()).Err()
	if err != nil {
		return err
	}

	return nil
}

func GetInstance() *redis.Client {
	if client == nil {
		once.Do(func() {
			if err := Initialize(); err != nil {
				panic("redis is not initialized")
			}
		})
	}

	return client
}
