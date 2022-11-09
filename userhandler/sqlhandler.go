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
	dbnum := ctx.PostForm("dbnum")
	dbname := ctx.PostForm("dbname")
	//log.Println(username, sql_str, dbname, dbnum)

	var (
		rowsInserted     int
		rowsUpdated      int
		rowsDeleted      int
		execResult       string
		globalRecoveryId string
		sqlRownum        int
		msg              string
		userid           int
		username         string
		err              error
	)
	username = ctx.GetString("username") //cros通过后上下文设置了
	userid = ctx.GetInt("userid")

	//此模块后期可以再加入JWT，传token，解析后再验证token中的用户
	userInWhitelist := uh.UserService.CheckUserInWhitelist(username)
	if !userInWhitelist {
		msg = "用户无权限"
		data := gin.H{
			"err": err,
		}
		response.Failed(ctx, data, msg)

		//写执行结果模块
		execResult = "Failed"
		InsertResultsErr := uh.UserService.Db.InsertResults(userid, username, execResult, msg, sqlRownum, rowsInserted, rowsUpdated, rowsDeleted, globalRecoveryId)
		if InsertResultsErr != nil {
			log.Println(InsertResultsErr)
		}

		return
	}
	//sql string长度限制，len>6
	if len(sql_str) <= 6 {
		msg := "please input correct SQL string"
		data := gin.H{}
		response.Failed(ctx, data, msg)
		ctx.Abort()
		return
	}

	//用户验证通过后，流程进入分析器sqlAnalyzer
	n, reason, isChecked := service.SqlAnalyzer(sql_str)
	if !isChecked {
		msg = "sql语句检测失败"
		data := gin.H{
			"reason":      reason,
			"sql -rownum": n,
		}
		response.Failed(ctx, data, msg)

		execResult = "Failed"
		InsertResultsErr := uh.UserService.Db.InsertResults(userid, username, execResult, msg, sqlRownum, rowsInserted, rowsUpdated, rowsDeleted, globalRecoveryId)
		if InsertResultsErr != nil {
			log.Println(InsertResultsErr)
		}

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
		if err != nil || scanRows > viper.GetInt("SqlExplainScanRowsLimit") {
			msg = "扫描检测失败"
			data := gin.H{
				"sql rownum": i,
				"scanrows":   scanRows,
				"error":      err,
			}
			response.Failed(ctx, data, msg)

			execResult = "Failed"
			sqlRownum = i
			InsertResultsErr := uh.UserService.Db.InsertResults(userid, username, execResult, msg, sqlRownum, rowsInserted, rowsUpdated, rowsDeleted, globalRecoveryId)
			if InsertResultsErr != nil {
				log.Println(InsertResultsErr)
			}

			ctx.Abort()
			return
		}
	}
	//检测阶段完全通过之后，进入备份阶段
	//生成全局recoveryid,只有前期检测全部通过
	randstr := utils.Randomstr(4, model.Letters)
	backupDir := viper.GetString("BackupDir")
	globalRecoveryId = "dms" + randstr + strconv.Itoa(int(time.Now().Unix()))

	for i := 1; i <= len(sqlmap); i++ {
		isValues := strings.Contains(sqlmap[i], "values")
		sqltype, _ := common.SqlTypeVerify(sqlmap[i])

		tag := "####第" + strconv.Itoa(i) + "条 " + sqltype + "#####" //写tag
		utils.FileWriter(globalRecoveryId, backupDir, tag)

		if isValues && sqltype == "insert" {
			//对于定值类insert的备份处理，直接把sql中value后的内容写入文件
			m := strings.IndexAny(sqlmap[i], "(")
			backupStr := sqlmap[i][m:]
			//log.Println(m, backupStr)
			err := utils.FileWriter(globalRecoveryId, backupDir, backupStr)
			if err != nil {
				msg = "执行前，写备份数据失败"
				data := gin.H{
					"err": err,
				}
				response.Failed(ctx, data, msg)

				sqlRownum = i
				execResult = "Failed"
				InsertResultsErr := uh.UserService.Db.InsertResults(userid, username, execResult, msg, sqlRownum, rowsInserted, rowsUpdated, rowsDeleted, globalRecoveryId)
				if InsertResultsErr != nil {
					log.Println(InsertResultsErr)
				}

				ctx.Abort()
				return
			}
		} else {
			newsql := common.SqlConvert2Select(sqlmap[i])
			//log.Println(newsql)
			resut, err := newUh.UserService.BackUpAndRecovery(newsql)
			//log.Println("result:", resut)
			if err == nil {
				for j := 0; j < len(resut); j++ {
					jsonResult := utils.Map2Json(resut[j])
					//log.Println("jsonResult:", jsonResult)
					err := utils.FileWriter(globalRecoveryId, backupDir, jsonResult)
					if err != nil {
						msg = "执行前，写备份数据失败"
						data := gin.H{
							"err": err,
						}
						response.Failed(ctx, data, msg)

						sqlRownum = i
						execResult = "Failed"
						InsertResultsErr := uh.UserService.Db.InsertResults(userid, username, execResult, msg, sqlRownum, rowsInserted, rowsUpdated, rowsDeleted, globalRecoveryId)
						if InsertResultsErr != nil {
							log.Println(InsertResultsErr)
						}

						ctx.Abort()
						return
					}
				}
			} else {
				msg = "执行前，获取备份数据失败"
				data := gin.H{
					"err": err,
				}
				response.Failed(ctx, data, msg)

				sqlRownum = i
				execResult = "Failed"
				InsertResultsErr := uh.UserService.Db.InsertResults(userid, username, execResult, msg, sqlRownum, rowsInserted, rowsUpdated, rowsDeleted, globalRecoveryId)
				if InsertResultsErr != nil {
					log.Println(InsertResultsErr)
				}
				ctx.Abort()
				return
			}

		}
	}
	/*检测备份通过之后，sql执行阶段；*/
	/*此处对于一系列传入的SQL并没有使用事物,视每条sql间没有事务依赖关系*/

	for i := 1; i <= len(sqlmap); i++ {
		resultRows, err := newUh.UserService.ExecSqlAndGetRownum(sqlmap[i])
		if err == nil {
			rowsUpdated = int(resultRows["updateRows"]) + rowsUpdated
			rowsDeleted = int(resultRows["deleteRows"]) + rowsDeleted
			rowsInserted = int(resultRows["insertRows"]) + rowsInserted
			log.Printf("###第%v条sql执行成功###%v插入行数:%v,更新行数:%v,删除行数:%v\n", i, sqlmap[i], int(resultRows["insertRows"]), int(resultRows["updateRows"]), int(resultRows["deleteRows"]))
		} else {
			//执行到任意行失败，则返回，并返回已经修改的行数，和错误信息
			msg := "执行中断"
			data := gin.H{
				"sql rownum":   i,
				"error":        err,
				"rowsInserted": rowsInserted,
				"rowsDeleted":  rowsDeleted,
				"rowsUpdated":  rowsUpdated,
			}
			log.Printf("###%v执行失败###\n", sqlmap[i])
			response.Failed(ctx, data, msg)

			sqlRownum = i
			execResult = "Failed"
			InsertResultsErr := uh.UserService.Db.InsertResults(userid, username, execResult, msg, sqlRownum, rowsInserted, rowsUpdated, rowsDeleted, globalRecoveryId)
			if InsertResultsErr != nil {
				log.Println(InsertResultsErr)
			}

			ctx.Abort()
			return
		}
	}
	//执行完成后在外层调用response.success统一返回
	msg = "执行成功"
	data := gin.H{
		"sql rownum":   len(sqlmap),
		"rowsInserted": rowsInserted,
		"rowsDeleted":  rowsDeleted,
		"rowsUpdated":  rowsUpdated,
	}
	response.Success(ctx, data, msg)

	sqlRownum = len(sqlmap)
	execResult = "Success"
	InsertResultsErr := uh.UserService.Db.InsertResults(userid, username, execResult, msg, sqlRownum, rowsInserted, rowsUpdated, rowsDeleted, globalRecoveryId)
	if InsertResultsErr != nil {
		log.Println(InsertResultsErr)
	}

	ctx.Abort()
}
