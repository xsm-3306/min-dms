package userhandler

import (
	"min-dms/common"
	"net/http"

	"github.com/gin-gonic/gin"
)

//sqlhandler 为所有SQL的总入口，所有的请求从此进入
func Sqlhandler(ctx *gin.Context) {
	//postform接收string
	sql_str := ctx.PostForm("sql")
	username := ctx.PostForm("username")

	//验证发送请求的用户，以白名单机制存放在库中
	isuserexists := common.Checkuserstatus(username)
	if !isuserexists {
		ctx.JSON(http.StatusNonAuthoritativeInfo, gin.H{
			"msg": "用户不在白名单内",
		})
		return
	}
	//对传进来的SQL进行一系列的验证

	ctx.JSON(http.StatusOK, gin.H{
		"code":     200,
		"sql":      sql_str,
		"username": username,
	})
}
