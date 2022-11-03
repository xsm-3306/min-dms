package dao

import (
	"errors"
	"fmt"
	"min-dms/common"
	"min-dms/model"
	"strconv"
)

//用户匹配,登录
//在dao层做password的比对
func (db *Database) CompareUserInfo(loguser *model.LoginUser) (token string, err error) {
	sql := "select id,username,password,is_deleted from user_info where username=?"
	var tokenUserInfo model.User

	//每一个err,都进入详细的流程，返回上层详细的err,有利于用户根据最上层的返回而做判断
	result, err := db.GetRows(sql, loguser.Username)
	if err != nil {
		return "", err
	}

	if len(result) < 1 {
		err = fmt.Errorf("user do not exists ")

		return "", err
	} else {
		if result[0]["is_deleted"] == "1" {
			return "", errors.New("user account is expired! ")
		}
		isPassword := common.PasswordVertify(loguser.Password, result[0]["password"])
		if isPassword {
			id, _ := strconv.Atoi(result[0]["id"])
			tokenUserInfo.Username = loguser.Username
			tokenUserInfo.Userid = id

			token, _ := common.GenToken(&tokenUserInfo)
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
	if err != nil {
		return err
	}

	return nil
}
