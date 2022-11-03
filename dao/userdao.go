package dao

import (
	"errors"
	"log"
	"min-dms/common"
	"strconv"
)

//检测是否白名单用户
func (db *Database) CheckUserInWhitelist(username string) bool {
	sql := "select id from user_whitelist where is_deleted=0 and username=?"

	result, err := db.GetRows(sql, username)
	//log.Println(result)
	if err == nil && len(result) > 0 {
		for key := range result[0] {
			if result[0][key] == "null_val" {
				//userid = 0
				return false
			} else {
				//userid, _ = strconv.Atoi(result[0][key])
				return true
			}
		}
	}
	return false
}

//检测用户表中用户是否存在
func (db *Database) CheckUserExists(username string) (bool, error) {
	sql := "select is_deleted from user_info where username=?"

	result, err := db.GetRows(sql, username)
	if err == nil && len(result) > 0 {
		for key := range result[0] {
			if result[0][key] == "0" {
				return true, errors.New("user account already exists! ")
			} else {
				return true, errors.New("user account already exists,and account is locked! ")
			}
		}
	}

	return false, err
}

//检查sql explain时扫描的行数
func (db *Database) CheckSqlExplainScanRows(sql string) (scanRows int, err error) {
	sqlstr := "explain " + sql
	scanRows = 1
	sqltype, _ := common.SqlTypeVerify(sql)

	result, err := db.GetRows(sqlstr)
	//log.Println(result)
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

//获取每个实例下的数据库列表
func (db *Database) GetDbList() (dbList []string, err error) {
	sql := "SELECT DISTINCT(table_schema) FROM information_schema.TABLES WHERE table_schema NOT IN('sys','system','mysql','information_schema','performance_schema')  ORDER BY table_schema"

	result, err := db.GetRows(sql)
	if err == nil {

		for i := 0; i < len(result); i++ {
			dbList = append(dbList, result[i]["table_schema"])
		}
		log.Println(dbList)
	}
	return
}

//执行结果入库
func (db *Database) InsertResults(vals ...interface{}) error {
	resutlInsertSql := "insert into user_sqlexec_log(user_id,username,exec_result,reason,sql_rownum,rows_inserted,rows_updated,rows_deleted,recovery_id)values(?,?,?,?,?,?,?,?,?)"
	_, err := db.AddRows(resutlInsertSql, vals...)
	if err != nil {
		return err
	}
	return nil
}
