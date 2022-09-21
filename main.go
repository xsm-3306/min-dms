package main

import (
	"min-dms/dao"
	userhandler "min-dms/userhandler"

	"github.com/gin-gonic/gin"
)

func main() {
	dao.InitDb()

	engine := gin.Default()
	engine.POST("/api/sqlhandler", func(ctx *gin.Context) {
		userhandler.SqlHandler(ctx)
	})
	engine.Run(":8081")
}
