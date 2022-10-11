package service

//备份本质上查库获取原数据并保存
func (us *UserService) BackUpAndRecovery(sql string) (result []map[string]string, err error) {
	return us.Db.GetRows(sql)
}
