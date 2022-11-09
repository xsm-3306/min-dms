package service

import "min-dms/model"

func (us *UserService) Login(loguser *model.LoginUser) (token string, err error) {
	return us.Db.CompareUserInfo(loguser)
}

func (us *UserService) Register(registeruser *model.LoginUser) error {
	return us.Db.AddUser(registeruser)
}

func (us *UserService) Logout(token string) (int, error) {
	return us.Db.UpdateTokenTab(token)
}
