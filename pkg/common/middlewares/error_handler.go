package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go-template/pkg/common/apierrors"
)

func ErrorHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
		if len(ctx.Errors) == 0 {
			return
		}

		err := ctx.Errors.Last().Err
		if errors.Is(err, apierrors.ErrorUnauthorized) {
			ctx.AbortWithStatusJSON(apierrors.UnauthorizedErrorCode, apierrors.ErrorUnauthorized)
			return
		}
		if errors.Is(err, apierrors.ErrorBadRequest) {
			ctx.AbortWithStatusJSON(apierrors.BadRequestErrorCode, apierrors.ErrorBadRequest)
			return
		}
		if errors.Is(err, apierrors.ErrorNotFound) {
			ctx.AbortWithStatusJSON(apierrors.NotFoundErrorCode, apierrors.ErrorNotFound)
			return
		}

		ctx.AbortWithStatusJSON(apierrors.InternalServerErrorCode, apierrors.ErrorInternalServer)
	}
}
