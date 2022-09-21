package utils

import (
	"strings"
)

//SplitStringByChar把字符串按需求分割成多个子串，并存入map[int]string中。
func SplitStringByChar(s string, substr string) map[int]string {
	s = strings.Trim(s, "\r\n")
	s = strings.Trim(s, "\n")
	s = strings.TrimSpace(s) //除去换行符及前后空格，linux及windows平台都考虑到

	n := strings.Count(s, substr)
	sqlmap := make(map[int]string, n)
	//sqlmap := model.SqlStatementMap
	for i := 1; i <= n; i++ {
		m := strings.Index(s, substr)
		subsql := s[0 : m+1]
		sqlmap[i] = subsql
		s = s[m+1:] //分割后重新赋值s
		s = strings.Trim(s, "\r\n")
		s = strings.Trim(s, "\n")
		s = strings.TrimSpace(s)
	}
	return sqlmap
}
