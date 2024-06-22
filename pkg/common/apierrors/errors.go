package apierrors

import (
	"github.com/pkg/errors"
)

var (
	ErrorBadRequest     = errors.New("bad request")
	ErrorUnauthorized   = errors.New("you are not allowed to do this operation")
	ErrorNotFound       = errors.New("not found")
	ErrorInternalServer = errors.New("internal server error")

	BadRequestErrorCode     = 400
	UnauthorizedErrorCode   = 401
	NotFoundErrorCode       = 404
	InternalServerErrorCode = 500
)

var ErrorCodes = map[error]int{
	ErrorUnauthorized:   UnauthorizedErrorCode,
	ErrorBadRequest:     BadRequestErrorCode,
	ErrorNotFound:       NotFoundErrorCode,
	ErrorInternalServer: InternalServerErrorCode,
}
