package model

//用来存放传递进来的多个或单个sql。
//遍历map时的元素顺序与添加键值对的顺序无关。
//所以对于业务意义上有先后顺序的SQL，必须保证遍历顺序
var SqlStatementMap map[int]string
