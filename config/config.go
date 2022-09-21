package config

type DbConfig struct {
	Host     string
	Port     int
	Db_name  string
	User     string
	Password string
}

//dms系统业业务数据库配置
var DmsConfig DbConfig
