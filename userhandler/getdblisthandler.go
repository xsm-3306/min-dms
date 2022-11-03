package userhandler

import (
	"log"
	"min-dms/dao"
	"min-dms/response"
	"min-dms/service"

	"github.com/gin-gonic/gin"
)

func (uh *Userhandler) GetDbList(ctx *gin.Context) {
	//username := ctx.PostForm("username")
	dbnum := ctx.PostForm("dbnum")

	username := ctx.GetString("username")

	log.Println(username, dbnum)
	//此模块后期可以再加入JWT，传token，解析后再验证token中的用户
	userInWhitelist := uh.UserService.CheckUserInWhitelist(username)
	if !userInWhitelist {
		msg := "用户无权限"
		data := gin.H{}
		response.Failed(ctx, data, msg)
		return
	}

	//根据传递数据，dbname dbnum,初始化新的数据库连接来执行相关的操作
	dbname := "mysql" //初始化需要dbname,给默认值，以建立连接
	newdb := new(dao.Database).NewDb(dbnum, dbname)
	newUs := service.UserService{Db: newdb}
	newUh := Userhandler{UserService: newUs}

	dbList, err := newUh.Db.GetDbList()
	if err == nil {
		msg := ""
		data := gin.H{
			"dbList": dbList,
		}
		response.Success(ctx, data, msg)
	} else {

		msg := "获取dbList失败"
		data := gin.H{
			"err": err,
		}
		response.Failed(ctx, data, msg)
		ctx.Abort()
		return
	}

	ctx.Abort()
}
