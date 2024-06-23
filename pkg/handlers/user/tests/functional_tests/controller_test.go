package functional_tests

import (
	"github.com/appleboy/gofight/v2"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

const (
	userNotFoundResponse = `{"error":"incorrect email or password: not found","success":false}`
)

func Test_Functional_Login_IncorrectEmail(t *testing.T) {
	gofight.
		New().
		POST("/api/user/login").
		SetJSON(gofight.D{
			"email":    "incorrect_email@test.app",
			"password": "123",
		}).
		Run(engine, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusNotFound, r.Code)
			assert.JSONEq(t, userNotFoundResponse, r.Body.String())
		})
}

func Test_Functional_Login_IncorrectPassword(t *testing.T) {
	gofight.
		New().
		POST("/api/user/login").
		SetJSON(gofight.D{
			"email":    "email@test.app",
			"password": "incorrect_password",
		}).
		Run(engine, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusNotFound, r.Code)
			assert.JSONEq(t, userNotFoundResponse, r.Body.String())
		})
}

func Test_Functional_Login_Success(t *testing.T) {
	gofight.
		New().
		POST("/api/user/login").
		SetJSON(gofight.D{
			"email":    "email@test.app",
			"password": "123",
		}).
		Run(engine, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, r.Code)
		})
}
