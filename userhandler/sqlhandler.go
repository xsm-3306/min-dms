package userhandler

import (
	"log"
	"min-dms/common"
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
	//对传进来的SQL进行一系列的验证

	n := common.SqlStatementSingleVerify(sql_str)
	if n > 5 {
		ctx.JSON(http.StatusNotAcceptable, gin.H{
			"msg": "不允许一次性超过10条SQL",
		})
		return
	}

	sql_type, err := common.SqlTypeVerify(sql_str)
	if err != nil {
		ctx.JSON(http.StatusNotAcceptable, gin.H{
			"msg": "非允许的SQL类型",
		})
		log.Println(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"sql":      sql_str,
		"username": username,
		"sql_type": sql_type,
	})
}
