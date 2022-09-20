package common

import (
	"errors"
	"min-dms/utils"
	"strings"
)

//Sqltypeverify 返回sql type，并且保证只能是delete update insert中的一种
func SqlTypeVerify(sql string) (string, error) {
	//去除头尾的空格，前六个字符即是sql type
	sql_str := strings.TrimSpace(sql)
	sql_type := sql_str[0:6]
	if sql_type == "insert" || sql_type == "delete" || sql_type == "update" {
		return sql_type, nil
	} else {
		err := errors.New("非允许的SQL类型")
		return "", err
	}
}

//Sqlstatementverify 返回sql语句的复杂程度；简单更新？或者是多表联合更新
func SqlStatementSimpleVerify() {

}

//SqlStatementSingleVerify 判断传入的SQL text为单条sql语句，或者是多条sql语句，返回sql条数N
func SqlStatementSingleVerify(sql string) int {
	sql_str := strings.TrimSpace(sql)
	last_char := sql_str[(len(sql_str) - 1):]
	if last_char != ";" {
		sql_str = sql_str + ";" //为SQL加结尾分号
	}
	sqlmap := utils.SplitStringByChar(sql_str, ";")

	return len(sqlmap)
}
