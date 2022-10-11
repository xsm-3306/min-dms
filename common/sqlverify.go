package common

import (
	"errors"
	"min-dms/model"
	"min-dms/utils"
	"strings"
)

//Sqltypeverify 检查sql type，确保只能是delete update insert中的一种，并且返回sqltype
func SqlTypeVerify(sql string) (sqlType string, err error) {

	//去除头尾的空格,默认前六个字符为sql_type(只能为insert,delete,update)
	sql_str := strings.TrimSpace(sql)
	sqlType = sql_str[0:6]
	if sqlType == "insert" || sqlType == "delete" || sqlType == "update" {
		//log.Println("验证的sql是：", sql)
		return sqlType, nil
	} else {
		err := errors.New("非允许的SQL类型")
		return "", err
	}
}

//Sqlstatementverify 返回sql语句的复杂程度；简单更新？或者是多表联合更新
//此处认为包含子查询，关联查询等都不能算简单SQL
//只有单表的简单查询，才符合要求
func SqlStatementSimpleVerify() {

}

//SqlLengthVerify 验证传入sql string的长度，最大允许长度定义再参数SqlLengthLimit中
func SqlLengthVerify(sql string) bool {

	return len(sql) < model.SqlLengthLimit

}

//
func SqlExplainScanRowsVerify(sql string) {

}

//SqlConvert2Select 实现SQL转换功能。
//insert delete update转换成对应where条件的select,在执行时的备份阶段会用到
func SqlConvert2Select(sql string) (newsql string) {

	sqltype, _ := SqlTypeVerify(sql)

	switch sqltype {
	case "delete":
		substr := utils.SplitStringByChar2(sql, "from")
		newsql = "select * from " + substr
	case "update":
		substr := utils.SplitStringByChar2(sql, "where")
		m := strings.Index(sql, "set")
		tablename := sql[6:m]
		newsql = "select * from" + tablename + " where " + substr

	case "insert":
		isValues := strings.Contains(sql, "values")
		if !isValues {
			m := strings.Index(sql, "select")
			newsql = sql[m:]
		}
	}

	return

}
