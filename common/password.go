package common

import (
	"fmt"
	"unicode"

	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

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

//密码强度校验规则
func PasswordStrengthVertify(pwdstr string) error {

	/*
		1.包含大写字母
		2.包含小写字母
		3.包含数字
		4.包含特殊字符
		5.长度在8-20之间，可通过配置文件配置
	*/
	var (
		isUpper   bool
		isLower   bool
		isNumber  bool
		isSpecial bool
	)

	pwdMaxLen := viper.GetInt("PasswordMaxLen")
	pwdMinLen := viper.GetInt("PasswordMinLen")

	if len(pwdstr) > pwdMaxLen || len(pwdstr) < pwdMinLen {
		return fmt.Errorf("the password lenght most bewteen %d and %d! ", pwdMinLen, pwdMaxLen)
	}

	//使用unicode包来对字符做判断
	for _, s := range pwdstr {
		switch {
		case unicode.IsUpper(s):
			isUpper = true
		case unicode.IsLower(s):
			isLower = true
		case unicode.IsNumber(s):
			isNumber = true
		case unicode.IsSymbol(s) || unicode.IsPunct(s):
			isSpecial = true
		default:
		}
	}

	if isUpper && isLower && isNumber && isSpecial {
		return nil
	}

	return fmt.Errorf("the password string most also contain upcase and lowcase letters,numbers and punctuation! ")
}
