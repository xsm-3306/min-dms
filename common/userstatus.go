package common

import (
	"log"
)

func Checkuserstatus(username string) bool {
	var userid int
	query_str := "select id from user_whitelist where is_deleted=0 and username=?"

	stmt, err := Con_pool.Prepare(query_str)
	if err != nil {
		log.Println("prepared failed", err)
		return false
	}
	defer stmt.Close()

	err = stmt.QueryRow(username).Scan(&userid)
	if err != nil {
		log.Println("query failed,no such user!", err)
		return false
	}

	return true
}
