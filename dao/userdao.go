package dao

import (
	"errors"
	"fmt"
	"log"
	"min-dms/common"
	"min-dms/model"
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

//用户匹配,登录
//在dao层做password的比对
func (db *Database) CompareUserInfo(loguser *model.LoginUser) (token string, err error) {
	sql := "select id,username,password,is_deleted from user_info where is_deleted=0 and username=?"

	//每一个err,都进入详细的流程，返回上层详细的err,有利于用户根据最上层的返回而做判断
	result, err := db.GetRows(sql, loguser.Username)
	if err != nil {
		return "", err
	}

	if len(result) < 1 {
		err = fmt.Errorf("user do not exists ")

		return "", err
	} else {
		isPassword := common.PasswordVertify(loguser.Password, result[0]["password"])
		if isPassword {
			token, _ := common.GenToken(loguser.Username)
			return token, nil
		} else {
			return "", errors.New("the password is not right")
		}
	}
}

//用户注册，账号密码及个人信息入库
func (db *Database) AddUser(registeruser *model.LoginUser) error {
	sql := "insert into user_info(username,password)values(?,?)"

	hashPassword, _ := common.PasswordHash(registeruser.Password)
	_, err := db.AddRows(sql, registeruser.Username, hashPassword)

	return err
}

//执行结果入库
func (db *Database) InsertResults(vals ...interface{}) error {
	resutlInsertSql := "insert into user_sqlexec_log(user_id,exec_result,reason,sql_rownum,rows_inserted,rows_updated,rows_deleted,recovery_id)values(?,?,?,?,?,?,?,?)"
	_, err := db.AddRows(resutlInsertSql, vals...)
	if err != nil {
		return err
	}
	return nil
}
