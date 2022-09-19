package main

import (
	common "min-dms/common"
	userhandler "min-dms/userhandler"

	"github.com/gin-gonic/gin"
)

func main() {
	common.InitDb()

	engine := gin.Default()
	engine.POST("/sqlhandler", func(ctx *gin.Context) {
		userhandler.Sqlhandler(ctx)
	})
	engine.Run(":8081")
}
