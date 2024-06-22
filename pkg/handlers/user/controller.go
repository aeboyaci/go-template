package user

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go-template/pkg/common/apierrors"
	"go-template/pkg/models"
	"net/http"
)

type Controller interface {
	Login(ctx *gin.Context)
	Register(ctx *gin.Context)
}

type controllerImpl struct {
	service Service
}

func NewController(service Service) Controller {
	return &controllerImpl{service}
}

func (c *controllerImpl) Login(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.Error(errors.Wrap(apierrors.ErrorBadRequest, err.Error()))
		return
	}

	token, err := c.service.Login(user.Email, user.Password)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"token":   token,
	})
}

func (c *controllerImpl) Register(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.Error(errors.Wrap(apierrors.ErrorBadRequest, err.Error()))
		return
	}

	if err := c.service.Register(user.Email, user.Password); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
	})
}
