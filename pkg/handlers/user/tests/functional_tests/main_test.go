package functional_tests

import (
	"github.com/gin-gonic/gin"
	"go-template/pkg/bootstrap"
	"go-template/pkg/common/utils"
	"os"
	"testing"
)

var engine *gin.Engine

func TestMain(m *testing.M) {
	err := bootstrap.Initialize()
	if err != nil {
		return
	}

	engine = gin.New()
	bootstrap.RegisterRoutes(engine)

	err = utils.LoadFixtures("./fixtures")
	if err != nil {
		return
	}

	os.Exit(m.Run())
}
