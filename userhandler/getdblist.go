package userhandler

import (
	"log"
	"min-dms/response"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

//获取已经配置了的可用的db list
func (uh *Userhandler) GetDbList(ctx *gin.Context) {
	username := ctx.PostForm("username")

	//此模块后期可以再加入JWT，传token，解析后再验证token中的用户
	userid, err := uh.UserService.GetUseridByUsername(username)
	if err != nil || userid < 1 {
		msg := "用户无权限"
		data := gin.H{}
		response.Failed(ctx, data, msg)
		return
	}

	//获取dblist,并返回给接口请求方
	var dblist []string
	err = viper.UnmarshalKey("dblist", &dblist)
	if err == nil {
		msg := ""
		data := gin.H{
			"dblist": dblist,
		}
		response.Success(ctx, data, msg)
	} else {
		log.Println("获取dblist失败，dblist解组失败", err)
		msg := "获取dblist失败"
		data := gin.H{
			"err": err,
		}
		response.Failed(ctx, data, msg)
		ctx.Abort()
		return
	}

	ctx.Abort()

}
