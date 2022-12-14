package engine

import (
	"min-dms/config"
	"min-dms/dao"
	"min-dms/middleware"
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
	dbNum := "db0" //本系统业务库放在db0,所以初始话的时候，先默认给志db0
	dbName := "dms"
	db.InitDbSource(dbNum, dbName)
	us = service.UserService{Db: db}
	uh = userhandler.Userhandler{UserService: us}

}

func InitEngine(engine *gin.Engine) {
	config.InitConfig()

	InitHandler()

	engine.Use(middleware.CrosMiddle())
	engine.POST("/api/user/login", uh.Login)
	engine.POST("/api/user/logout", uh.Logout)
	engine.POST("/api/user/register", uh.Register)

	router := engine.Group("/")
	router.Use(middleware.JwtAuthMiddle(&uh))
	{
		router.POST("/api/sqlhandler", uh.SqlHandler)
		router.POST("/api/getdbinstancelist", uh.GetDbInstanceList)
		router.POST("/api/getdblist", uh.GetDbList)
	}

}
