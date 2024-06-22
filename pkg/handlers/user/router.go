package user

import (
	"github.com/gin-gonic/gin"
	"go-template/pkg/common/database"
	"go-template/pkg/common/logger"
)

func Register(apiGroup *gin.RouterGroup) {
	controller := NewController(
		NewService(
			database.GetInstance(),
			logger.GetInstance(),
			NewRepository(),
		),
	)

	userGroup := apiGroup.Group("/user")
	userGroup.POST("/login", controller.Login)
	userGroup.POST("/register", controller.Register)
}
