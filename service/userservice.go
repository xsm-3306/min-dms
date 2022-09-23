package service

import "min-dms/dao"

type UserService struct {
	Db *dao.Database
}

func (us *UserService) GetUseridByUsername(username string) (userid int, err error) {
	return us.Db.GetUseridByUsername(username)
}
