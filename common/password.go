package common

import "golang.org/x/crypto/bcrypt"

/*bcrypt对数据密码加解密*/

//密码加密
func PasswordHash(pwdstr string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwdstr), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

//密码验证
func PasswordVertify(pwdstr string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwdstr))
	return err == nil
}
