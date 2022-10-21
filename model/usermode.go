package model

//user struct
type User struct {
	Userid   int
	Username string
}

//login struct bind with json and form
type LoginUser struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

//生成随机字符串的种子
var Letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
