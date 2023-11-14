package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rekib0023/event-horizon-gateway/controller"
)

func main() {
	if os.Getenv("ENVIRONMENT") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	controller.Start()
}
