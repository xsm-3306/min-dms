package engine

import (
	"min-dms/config"
	"min-dms/dao"
	"min-dms/service"
	"min-dms/userhandler"

	"github.com/gin-gonic/gin"
)

var (
	us service.UserService
	uh userhandler.Userhandler
)

func InitHandler() {
	db := &dao.Database{}
	db.InitDbSource()
	us = service.UserService{Db: db}
	uh = userhandler.Userhandler{UserService: us}

}

func InitEngine(engine *gin.Engine) {
	config.InitConfig()

	InitHandler()

	engine.POST("/api/sqlhandler", uh.SqlHandler)
}
