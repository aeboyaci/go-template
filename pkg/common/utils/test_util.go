package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/go-testfixtures/testfixtures/v3"
	"go-template/pkg/common/database"
	"go-template/pkg/common/logger"
	"go.uber.org/zap"
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

func LoadFixtures(fixturesDirectory string) error {
	db, err := database.GetInstance().DB()
	if err != nil {
		logger.GetInstance().Error("Failed to get database instance", zap.String("location", "LoadFixtures"), zap.Error(err))
		return err
	}

	fixtures, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Directory(fixturesDirectory),
	)
	if err != nil {
		logger.GetInstance().Error("Failed to create fixtures instance", zap.String("location", "LoadFixtures"), zap.Error(err))
		return err
	}

	err = fixtures.EnsureTestDatabase()
	if err != nil {
		logger.GetInstance().Error("Failed to ensure test database", zap.String("location", "LoadFixtures"), zap.Error(err))
		return err
	}

	err = fixtures.Load()
	if err != nil {
		logger.GetInstance().Error("Failed to load fixtures", zap.String("location", "LoadFixtures"), zap.Error(err))
		return err
	}

	return nil
}
