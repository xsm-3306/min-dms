package utils

import (
	"encoding/json"
	"log"
	"os"
	"runtime"
	"strings"
)

//SplitStringByChar把字符串按需求分割成多个子串，并存入map[int]string中。
func SplitStringByChar(s string, substr string) map[int]string {
	s = strings.ToLower(s)
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

//写文件，根据路径以及文件名,每次追加写，文件不存在则创建,
func FileWriter(filename string, path string, writestr string) error {
	filePath := path + filename
	switch runtime.GOOS { //根据操作系系统类型添加换行符
	case "windows":
		writestr = writestr + "\n"
	default:
		writestr = writestr + "\r\n"
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Println("创建文件失败：", filePath)
		return err
	}
	defer file.Close()

	_, err = file.WriteString(writestr)
	if err != nil {
		return err
	}
	file.Sync()
	return nil
}

//map2json,map[string][string]转换为json；dao.GetRows返回的是map[string]string,写到文件中的可读性不高
func Map2Json(result map[string]string) (jsonstr string) {
	js, _ := json.Marshal(result)
	jsonstr = string(js)
	return
}
