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
		statusCode, ok := apierrors.ErrorCodes[errors.Cause(err)]
		if !ok {
			statusCode = apierrors.InternalServerErrorCode
		}

		ctx.AbortWithStatusJSON(statusCode, gin.H{
			"success": false,
			"error":   err.Error(),
		})
	}
}
