package dao

import "strconv"

func (db *Database) GetUseridByUsername(username string) (userid int, err error) {
	sql := "select id from user_whitelist where is_deleted=0 and username=?"

	result, err := db.GetRows(sql, username)
	if err == nil {
		for key := range result[0] {
			if result[0][key] == "null_val" {
				userid = 0
			} else {
				userid, _ = strconv.Atoi(result[0][key])
			}
		}
	}
	return
}
