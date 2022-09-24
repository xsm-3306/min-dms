package service

import (
	"min-dms/common"
)

//sql分析器，sql分析阶段
func SqlAnalyzer(sql string) (int, string, bool) {
	//先拆分为单条sql，存放再sqlMap[int]string中
	sqlMap := common.SqlStatementSplit(sql)

	for i := 1; i <= len(sqlMap); i++ {
		//log.Println(i, sqlMap[i])
		//验证sql type类型
		err := common.SqlTypeVerify(sqlMap[i])
		if err != nil {
			reason := "不允许的sql类型"
			return i, reason, false
		}
		//验证长度
		isChecked := common.SqlLengthVerify(sqlMap[i])
		if !isChecked {
			reason := "sql长度超过了允许范围"
			return i, reason, false
		}

	}

	return 0, "", true
}
