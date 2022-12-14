package userhandler

import (
	"min-dms/common"
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
	//log.Println(loginuser)

	token, err := uh.UserService.Login(&loginuser)
	if err != nil {
		data := gin.H{
			"err": err.Error(),
		}
		msg := "用户验证失败"
		response.Failed(ctx, data, msg)
		ctx.Abort()
		return
	}
	//获取到的token入库，后续可以做成类似token注销的功能；更建议直接存在redis中。
	//但是此处不影响主流程，即即使没有写入，也不影响后续
	sqlstr := "insert into user_authtoken_log(token_str)values(?)"
	uh.UserService.Db.AddRows(sqlstr, token)

	data := gin.H{
		"token": token,
	}
	msg := "登录成功"
	response.Success(ctx, data, msg)

	ctx.Abort()
}

//登出
func (uh *Userhandler) Logout(ctx *gin.Context) {
	token := ctx.Request.Header.Get("Authorization")

	if len(token) <= 7 {
		msg := "the auth header is empty"
		data := gin.H{}
		response.Failed(ctx, data, msg)
		ctx.Abort()
		return
	}
	token = token[7:]
	rowsUpdated, err := uh.UserService.Logout(token)
	if err != nil {
		data := gin.H{
			"err": err,
		}
		msg := "登出失败"
		response.Failed(ctx, data, msg)
		ctx.Abort()
		return
	} else {
		if rowsUpdated == 0 {
			data := gin.H{
				"rowsUpdated": rowsUpdated,
			}
			msg := "登出失败,不存在登录状态"
			response.Failed(ctx, data, msg)
			ctx.Abort()
			return
		} else {
			data := gin.H{
				"rowsUpdated": rowsUpdated,
			}
			msg := "登出成功"
			response.Success(ctx, data, msg)
			ctx.Abort()
			return
		}
	}
}

//用户注册handler
func (uh *Userhandler) Register(ctx *gin.Context) {
	var registeruser model.LoginUser
	//解析
	if err := ctx.ShouldBind(&registeruser); err != nil {
		data := gin.H{
			"err": err,
		}
		msg := "解析账号密码错误"

		response.Failed(ctx, data, msg)
		ctx.Abort()
		return
	}
	//先检测用户是否存在
	userExists, err := uh.UserService.Db.CheckUserExists(registeruser.Username)
	if userExists {
		data := gin.H{
			"err": err.Error(),
		}
		msg := "用户账号已存在"
		response.Failed(ctx, data, msg)
		ctx.Abort()
		return
	}

	//验证
	err = common.PasswordStrengthVertify(registeruser.Password)
	if err != nil {
		data := gin.H{
			"err": err.Error(),
		}
		msg := "密码强度不符合要求"
		response.Failed(ctx, data, msg)
		ctx.Abort()
		return
	}
	//注册入库
	err = uh.UserService.Db.AddUser(&registeruser)
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
