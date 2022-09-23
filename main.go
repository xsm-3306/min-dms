package main

import (
	"min-dms/engine"

	"github.com/gin-gonic/gin"
)

func main() {

	dms := gin.Default()
	engine.InitEngine(dms)
	dms.Run(":8081")
}
