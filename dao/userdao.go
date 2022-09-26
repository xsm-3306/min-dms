package dao

import (
	"log"
	"min-dms/common"
	"strconv"
)

//根据用户名查询id
func (db *Database) GetUseridByUsername(username string) (userid int, err error) {
	sql := "select id from user_whitelist where is_deleted=0 and username=?"

	result, err := db.GetRows(sql, username)
	//log.Println(result)
	if err == nil && len(result) > 0 {
		for key := range result[0] {
			if result[0][key] == "null_val" {
				userid = 0
			} else {
				userid, _ = strconv.Atoi(result[0][key])
			}
		}
	}
	return
}

//检查sql explain时扫描的行数
func (db *Database) CheckSqlExplainScanRows(sql string) (scanRows int, err error) {
	sqlstr := "explain " + sql
	scanRows = 1
	sqltype, _ := common.SqlTypeVerify(sql)

	result, err := db.GetRows(sqlstr)
	log.Println(sqlstr)

	if err == nil { //explain结果根据sqltype分别处理
		switch sqltype {
		case "insert": //insert 类型对于null_val需要单独处理
			for i := 0; i < len(result); i++ {
				if result[i]["rows"] == "null_val" {
					rows := 1
					scanRows = scanRows * rows
				} else {
					rows, err := strconv.Atoi(result[i]["rows"]) //取扫描行数
					if err == nil {
						scanRows = scanRows * rows
					} else {
						return 0, err
					}
				}
			}
		case "delete":
			for i := 0; i < len(result); i++ {
				rows, err := strconv.Atoi(result[i]["rows"]) //取扫描行数
				if err == nil {
					scanRows = scanRows * rows
				} else {
					return 0, err
				}
			}
		case "updaet":
			for i := 0; i < len(result); i++ {
				rows, err := strconv.Atoi(result[i]["rows"]) //取扫描行数
				if err == nil {
					scanRows = scanRows * rows
				} else {
					return 0, err
				}
			}
		}

	}
	return
}
