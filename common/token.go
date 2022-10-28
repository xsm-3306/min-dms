package common

import (
	"errors"
	"min-dms/model"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type UserClaims struct {
	model.User
	jwt.RegisteredClaims
}

//生成token
func GenToken(username string) (string, error) {

	claims := UserClaims{
		User: model.User{
			Username: username,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "iamchaos",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(model.TokenExpiredTime)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(model.SignStr)
}

//解析token
func ParseToken(tokenstr string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenstr, &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
		return model.SignStr, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("parse token failed")
}
