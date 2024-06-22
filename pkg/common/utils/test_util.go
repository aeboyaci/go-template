package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
)

func SetupTestGinContext(req *http.Request) *gin.Context {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	return c
}
