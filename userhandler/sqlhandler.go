package userhandler

import (
	"min-dms/common"
	"min-dms/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

//sqlhandler 为所有SQL的总入口，所有的请求从此进入
func SqlHandler(ctx *gin.Context) {
	//postform接收string
	sql_str := ctx.PostForm("sql")
	username := ctx.PostForm("username")

	//此模块后期可以用jwt token的方式替代，传token，解析后再验证token中的用户
	isuserexists := common.CheckUserStatus(username)
	if !isuserexists {
		ctx.JSON(http.StatusNotAcceptable, gin.H{
			"msg": "用户不在白名单内",
		})
		return
	}

	//用户验证通过后，流程进入分析器sqlAnalyzer
	n, reason, isChecked := service.SqlAnalyzer(sql_str)
	if !isChecked {
		ctx.JSON(http.StatusMethodNotAllowed, gin.H{
			"msg":    "检测不通过",
			"位置":     n,
			"reason": reason,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"sql":      sql_str,
		"username": username,
	})
}
