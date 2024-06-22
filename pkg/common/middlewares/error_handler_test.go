package middlewares

import (
	"github.com/stretchr/testify/assert"
	"go-template/pkg/common/apierrors"
	"go-template/pkg/common/utils"
	"net/http"
	"testing"
)

func Test_ErrorHandler_Unauthorized(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	c := utils.SetupTestGinContext(req)
	c.Error(apierrors.ErrorUnauthorized)

	ErrorHandler()(c)

	assert.Equal(t, c.Writer.Status(), http.StatusUnauthorized)
}

func Test_ErrorHandler_BadRequest(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	c := utils.SetupTestGinContext(req)
	c.Error(apierrors.ErrorBadRequest)

	ErrorHandler()(c)

	assert.Equal(t, c.Writer.Status(), http.StatusBadRequest)
}

func Test_ErrorHandler_NotFound(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	c := utils.SetupTestGinContext(req)
	c.Error(apierrors.ErrorNotFound)

	ErrorHandler()(c)

	assert.Equal(t, c.Writer.Status(), http.StatusNotFound)
}

func Test_ErrorHandler_InternalServerError(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	c := utils.SetupTestGinContext(req)
	c.Error(apierrors.ErrorInternalServer)

	ErrorHandler()(c)

	assert.Equal(t, c.Writer.Status(), http.StatusInternalServerError)
}

func Test_ErrorHandler_Success(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	c := utils.SetupTestGinContext(req)

	ErrorHandler()(c)

	assert.Equal(t, c.Writer.Status(), http.StatusOK)
}
