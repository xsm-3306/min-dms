package common

import (
	"min-dms/utils"
	"strings"
)

//SqlStatementSplit 拆分传入的sql string,拆分单条后存入sqlMap[int]string
//map底层的hashmap做hash散列时并没有保证有序，所以用range遍历的时注意业务是否有顺序的需要
func SqlStatementSplit(sql string) map[int]string {
	sql_str := strings.TrimSpace(sql)
	last_char := sql_str[(len(sql_str) - 1):]
	if last_char != ";" {
		sql_str = sql_str + ";" //为SQL加结尾分号
	}
	n := strings.Count(sql_str, ";")
	sqlmap := utils.SplitStringByChar(sql_str, ";", n)
	return sqlmap
}
