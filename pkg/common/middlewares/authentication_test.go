package middlewares

import (
	"github.com/stretchr/testify/assert"
	"go-template/pkg/common/apierrors"
	"go-template/pkg/common/redis"
	"go-template/pkg/common/utils"
	"go-template/pkg/handlers/user"
	"go-template/pkg/models"
	"gorm.io/gorm"
	"net/http"
	"testing"
	"time"
)

func Test_EnforceAuthentication_InvalidHeader(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	c := utils.SetupTestGinContext(req)
	EnforceAuthentication(nil, redis.GetInstance(), user.NewMockRepository())(c)

	assert.Equal(t, c.Errors.Last().Err, apierrors.ErrorUnauthorized)
}

func Test_EnforceAuthentication_InvalidToken(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header = http.Header{"Authorization": []string{"invalid_token"}}
	c := utils.SetupTestGinContext(req)
	EnforceAuthentication(nil, redis.GetInstance(), user.NewMockRepository())(c)

	assert.Equal(t, c.Errors.Last().Err, apierrors.ErrorUnauthorized)
}

func Test_EnforceAuthentication_ExpiredToken(t *testing.T) {
	token := utils.SignToken(1, time.Now().Add(-1*time.Hour))

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header = http.Header{"Authorization": []string{token}}
	c := utils.SetupTestGinContext(req)
	EnforceAuthentication(nil, redis.GetInstance(), user.NewMockRepository())(c)

	assert.Equal(t, c.Errors.Last().Err, apierrors.ErrorUnauthorized)
}

func Test_EnforceAuthentication_UserNotFound(t *testing.T) {
	token := utils.SignToken(1, time.Now().Add(1*time.Hour))

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header = http.Header{"Authorization": []string{token}}
	c := utils.SetupTestGinContext(req)

	mockUserRepository := user.NewMockRepository()
	mockUserRepository.MFindUserById = func(tx *gorm.DB, userId uint) (models.User, error) {
		return models.User{}, gorm.ErrRecordNotFound
	}
	EnforceAuthentication(nil, redis.GetInstance(), mockUserRepository)(c)

	assert.Equal(t, c.Errors.Last().Err, apierrors.ErrorUnauthorized)
}

func Test_EnforceAuthentication_ValidToken(t *testing.T) {
	token := utils.SignToken(1, time.Now().Add(1*time.Hour))

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header = http.Header{"Authorization": []string{token}}
	c := utils.SetupTestGinContext(req)

	mockUserRepository := user.NewMockRepository()
	mockUserRepository.MFindUserById = func(tx *gorm.DB, userId uint) (models.User, error) {
		return models.User{}, nil
	}
	EnforceAuthentication(nil, redis.GetInstance(), mockUserRepository)(c)

	assert.Equal(t, len(c.Errors), 0)
}
