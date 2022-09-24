package userhandler

import (
	"min-dms/common"
	"min-dms/model"
	"min-dms/response"
	"min-dms/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Userhandler struct {
	service.UserService
}

//sqlhandler 为所有SQL的总入口，所有的请求从此进入
func (uh *Userhandler) SqlHandler(ctx *gin.Context) {
	//postform接收string
	sql_str := ctx.PostForm("sql")
	username := ctx.PostForm("username")

	//此模块后期可以再加入JWT，传token，解析后再验证token中的用户
	userid, err := uh.UserService.GetUseridByUsername(username)
	if err != nil || userid < 1 {
		msg := "用户无权限"
		data := gin.H{}
		response.Failed(ctx, data, msg)
		return
	}

	//用户验证通过后，流程进入分析器sqlAnalyzer
	n, reason, isChecked := service.SqlAnalyzer(sql_str)
	if !isChecked {
		ctx.JSON(http.StatusMethodNotAllowed, gin.H{
			"msg":    "检测不通过",
			"第几行":    n,
			"reason": reason,
		})
		return
	}

	//explain 扫描行数检测；先拆分，再逐一检测，任何一个不符合规定，返回
	sqlmap := common.SqlStatementSplit(sql_str)
	for i := 1; i <= len(sqlmap); i++ {
		scanRows, err := uh.UserService.CheckSqlExplainScanRows(sqlmap[i])
		if err != nil || scanRows > model.SqlExplainScanRowsLimit {
			msg := "扫描行数检测失败"
			data := gin.H{
				"位置,第几行": i,
				"扫描行数":   scanRows,
				"error":  err,
			}
			response.Failed(ctx, data, msg)
			ctx.Abort()
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{

		"username": username,
	})
}
