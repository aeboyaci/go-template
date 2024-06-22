package bootstrap

import (
	"github.com/gin-gonic/gin"
	"go-template/pkg/common/database"
	"go-template/pkg/common/env"
	"go-template/pkg/common/logger"
	"go-template/pkg/common/middlewares"
	"go-template/pkg/common/redis"
	"go-template/pkg/handlers/user"
	"go.uber.org/zap"
)

func Initialize() error {
	var err error
	err = logger.Initialize()
	if err != nil {
		return err
	}

	err = env.Load()
	if err != nil {
		logger.GetInstance().Error("Failed to load environment variables", zap.Error(err))
		return err
	}

	err = database.Initialize()
	if err != nil {
		logger.GetInstance().Error("Failed to initialize database", zap.Error(err))
		return err
	}

	err = redis.Initialize()
	if err != nil {
		logger.GetInstance().Error("Failed to initialize Redis", zap.Error(err))
		return err
	}

	logger.GetInstance().Info("Application initialized successfully")
	return nil
}

func RegisterRoutes(engine *gin.Engine) {
	engine.Use(middlewares.ErrorHandler())

	apiGroup := engine.Group("/api")

	protectedApiGroup := engine.Group("/api")
	protectedApiGroup.Use(middlewares.EnforceAuthentication(redis.GetInstance(), user.NewRepository()))

	user.Register(apiGroup)
}
