package model

//用来存放传递进来的多个或单个sql。
//遍历map时的元素顺序与添加键值对的顺序无关。
//所以对于业务意义上有先后顺序的SQL，必须保证遍历顺序
var SqlStatementMap map[int]string

//一次性允许执行的SQL条数限制
var SqlRowsLimit = 50

//单条SQL允许的长度，以len()函数的结果为准，取的是字节数，
//所以SQL中含有中文的时候注意识别
var SqlLengthLimit = 1000

//限制每条sql explain时预估扫描的行数
//对于多条sql同时执行的情况，模式可以稍作调整
var SqlExplainScanRowsLimit = 10000

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
