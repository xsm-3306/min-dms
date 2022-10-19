package userhandler

import (
	"log"
	"min-dms/model"
	"min-dms/response"

	"github.com/gin-gonic/gin"
)

//login，验证账户可用性，密码正确性，并返回token，后续步骤使用此token进行访问
func (uh *Userhandler) Login(ctx *gin.Context) {
	var loginuser model.LoginUser
	//bind loginfo
	if err := ctx.ShouldBind(&loginuser); err != nil {
		data := gin.H{
			"err": err,
		}
		msg := "解析账号密码错误"

		response.Failed(ctx, data, msg)
		ctx.Abort()
		return
	}
	log.Println(loginuser)

	token, err := uh.UserService.Login(&loginuser)
	if err != nil {
		data := gin.H{
			"err": err,
		}
		msg := "用户验证失败"
		response.Failed(ctx, data, msg)
		ctx.Abort()
		return
	}

	data := gin.H{
		"token": token,
	}
	msg := ""
	response.Success(ctx, data, msg)

	ctx.Abort()
}

//用户注册handler
func (uh *Userhandler) Register(ctx *gin.Context) {
	var registeruser model.LoginUser

	if err := ctx.ShouldBind(&registeruser); err != nil {
		data := gin.H{
			"err": err,
		}
		msg := "解析账号密码错误"

		response.Failed(ctx, data, msg)
		ctx.Abort()
		return
	}

	err := uh.UserService.Db.AddUser(&registeruser)
	if err != nil {
		data := gin.H{
			"err": err,
		}
		msg := "注册失败"
		response.Failed(ctx, data, msg)
		ctx.Abort()
		return
	}

	data := gin.H{}
	msg := "注册成功"
	response.Success(ctx, data, msg)

	ctx.Abort()
}
