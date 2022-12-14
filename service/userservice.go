package service

import "min-dms/dao"

type UserService struct {
	Db *dao.Database
}

func (us *UserService) CheckUserInWhitelist(username string) bool {
	return us.Db.CheckUserInWhitelist(username)
}

func (us *UserService) CheckSqlExplainScanRows(sql string) (scanRows int, err error) {
	return us.Db.CheckSqlExplainScanRows(sql)
}

func (us *UserService) GetDbList() (dbList []string, err error) {
	return us.Db.GetDbList()
}
