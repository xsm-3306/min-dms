package dao

import (
	"errors"
	"fmt"
	"min-dms/common"
	"min-dms/model"
)

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
