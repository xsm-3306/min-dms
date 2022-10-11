package userhandler

import (
	"log"
	"min-dms/response"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

//获取已经配置了的可用的db instance list,并不是查库；而是从初始化的配置文件中读取并返回
func (uh *Userhandler) GetDbInstanceList(ctx *gin.Context) {
	username := ctx.PostForm("username")
	log.Println(username)
	//此模块后期可以再加入JWT，传token，解析后再验证token中的用户
	userid, err := uh.UserService.GetUseridByUsername(username)
	if err != nil || userid < 1 {
		msg := "用户无权限"
		data := gin.H{}
		response.Failed(ctx, data, msg)
		return
	}

	//获取dblist,并返回给接口请求方
	var dbNumList []string
	err = viper.UnmarshalKey("dblist", &dbNumList)
	if err == nil {
		msg := ""
		data := gin.H{
			"dbNumList": dbNumList,
		}
		response.Success(ctx, data, msg)
	} else {
		log.Println("获取dbNumList失败，dbNumList解组失败", err)
		msg := "获取dbNumList失败"
		data := gin.H{
			"err": err,
		}
		response.Failed(ctx, data, msg)
		ctx.Abort()
		return
	}

	ctx.Abort()

}
