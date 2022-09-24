package dao

import (
	"log"
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
	sql_str := "explain " + sql
	scanRows = 1

	result, err := db.GetRows(sql_str)
	log.Println(sql_str)
	if err == nil {

		for i := 0; i < len(result); i++ {
			//log.Println(len(result), scanRows)
			//log.Println(result[i])
			rows, err := strconv.Atoi(result[i]["rows"]) //取结果集中rows的值
			if err == nil {
				scanRows = scanRows * rows
			} else {
				return 0, err
			}

		}
	}
	return
}
