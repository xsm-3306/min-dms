package service

//
func (us *UserService) ExecSqlAndGetRownum(sql string) (resultRows map[string]int64, err2 error, err1 error) {
	return us.Db.ExecSql(sql)
}
