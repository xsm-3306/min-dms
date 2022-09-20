package utils

import (
	"strings"
)

//SplitStringByChar把字符串按需求分割成多个子串，并存入map中。
//对于有业务意义的情况，注意map的遍历顺序
func SplitStringByChar(s string, substr string) map[int]string {
	n := strings.Count(s, substr)
	sqlmap := make(map[int]string, n)
	//sqlmap := model.SqlStatementMap
	for i := 1; i <= n; i++ {
		m := strings.Index(s, substr)
		subsql := s[0 : m+1]
		sqlmap[i] = subsql
		s = s[m+1:]
	}
	return sqlmap
}
