package main

import (
	"github.com/gin-gonic/gin"
	"go-template/pkg/bootstrap"
)

func main() {
	err := bootstrap.Initialize()
	if err != nil {
		return
	}

	engine := gin.New()
	bootstrap.RegisterRoutes(engine)
	engine.Run(":8080")
}
