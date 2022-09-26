package service

//
func (us *UserService) ExecSqlAndGetRownum(sql string) (resultRows map[string]int64, err error) {
	return us.Db.ExecSql(sql)
}
