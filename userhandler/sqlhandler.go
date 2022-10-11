package userhandler

import (
	"log"
	"min-dms/common"
	"min-dms/dao"
	"min-dms/model"
	"min-dms/response"
	"min-dms/service"
	"min-dms/utils"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type Userhandler struct {
	service.UserService
}

//sqlhandler 为所有SQL的总入口，所有的请求从此进入
func (uh *Userhandler) SqlHandler(ctx *gin.Context) {
	//postform接收string
	sql_str := ctx.PostForm("sql")
	username := ctx.PostForm("username")
	dbnum := ctx.PostForm("dbnum")
	dbname := ctx.PostForm("dbname")
	log.Println(username, sql_str, dbname, dbnum)

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
		msg := "sql语句检测失败"
		data := gin.H{
			"reason": reason,
			"rownum": n,
		}
		response.Failed(ctx, data, msg)
		ctx.Abort()
		return
	}

	//根据传递数据，dbname dbnum,初始化新的数据库连接来执行相关的操作
	newdb := new(dao.Database).NewDb(dbnum, dbname)
	newUs := service.UserService{Db: newdb}
	newUh := Userhandler{UserService: newUs}

	//explain 扫描行数检测；先拆分，再逐一检测，任何一个不符合规定，返回
	sqlmap := common.SqlStatementSplit(sql_str)
	for i := 1; i <= len(sqlmap); i++ {

		scanRows, err := newUh.UserService.CheckSqlExplainScanRows(sqlmap[i])
		if err != nil || scanRows > model.SqlExplainScanRowsLimit {
			msg := "扫描检测失败"
			data := gin.H{
				"rownum":   i,
				"scanrows": scanRows,
				"error":    err,
			}
			response.Failed(ctx, data, msg)
			ctx.Abort()
			return
		}
	}
	//检测阶段完全通过之后，进入备份阶段
	randstr := utils.Randomstr(4, model.Letters)
	backupDir := viper.GetString("BackupDir")
	GlobalRecoveryId := "dms" + randstr + strconv.Itoa(int(time.Now().Unix()))

	for i := 1; i <= len(sqlmap); i++ {
		sqltype, _ := common.SqlTypeVerify(sqlmap[i])
		if sqltype == "insert" {
			isValues := strings.Contains(sqlmap[i], "values")
			if isValues {
				utils.FileWriter(GlobalRecoveryId, backupDir, sqlmap[i])
			}
		}
	}

	/*检测备份通过之后，sql执行阶段；*/
	/*此处对于一系列传入的SQL并没有使用事物,视每条sql间没有事务依赖关系*/
	var (
		rowsInserted int
		rowsUpdated  int
		rowsDeleted  int
	)
	for i := 1; i <= len(sqlmap); i++ {
		resultRows, err := newUh.UserService.ExecSqlAndGetRownum(sqlmap[i])
		if err == nil {
			rowsUpdated = int(resultRows["updateRows"]) + rowsUpdated
			rowsDeleted = int(resultRows["deleteRows"]) + rowsDeleted
			rowsInserted = int(resultRows["insertRows"]) + rowsInserted
			log.Printf("###第%v条sql执行成功###%v插入行数:%v,更新行数:%v,删除行数:%v\n", i, sqlmap[i], rowsInserted, rowsUpdated, rowsDeleted)
		} else {
			//执行到任意行失败，则返回，并返回已经修改的行数，和错误信息
			msg := "执行中断"
			data := gin.H{
				"rownum":       i,
				"error":        err,
				"rowsInserted": rowsInserted,
				"rowsDeleted":  rowsDeleted,
				"rowsUpdated":  rowsUpdated,
			}
			log.Printf("###%v执行失败###\n", sqlmap[i])
			response.Failed(ctx, data, msg)
			ctx.Abort()
			return
		}
	}
	//执行完成后在外层调用response.success统一返回
	msg := "执行成功"
	data := gin.H{
		"rownum":       len(sqlmap),
		"rowsInserted": rowsInserted,
		"rowsDeleted":  rowsDeleted,
		"rowsUpdated":  rowsUpdated,
	}
	response.Success(ctx, data, msg)

	ctx.Abort()
}
