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
		msg := "在执行前的分析阶段，sql检测不通过"
		data := gin.H{
			"reason": reason,
			"第几行":    n,
		}
		response.Failed(ctx, data, msg)
		ctx.Abort()
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

	/*检测通过之后，sql执行阶段；*/
	/*此处对于一系列传入的SQL并没有使用事物*/
	for i := 1; i <= len(sqlmap); i++ {
		resultRows, err1, err2 := uh.UserService.ExecSqlAndGetRownum(sqlmap[i])
		if err1 != nil {
			msg := "执行失败"
			data := gin.H{
				"sql位置,第几行":  i,
				"error1":     err1,
				"resultRows": resultRows,
				"error2":     err2,
			}
			response.Failed(ctx, data, msg)
			ctx.Abort()
			return
		} else {
			msg := "执行成功"
			data := gin.H{
				"sql总数":      i,
				"error1":     err1,
				"resultRows": resultRows,
				"error2":     err2,
			}
			response.Success(ctx, data, msg)
		}
	}

	ctx.JSON(http.StatusOK, gin.H{

		"username": username,
	})
	ctx.Abort()
}
