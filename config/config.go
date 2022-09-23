package config

//dms系统db配置
var (
	DmsDbUrl      string
	DbName        string
	DmsDbUser     string
	DmsDbPassword string
)

func InitConfig() {
	DmsDbUrl = "192.168.19.39:3306"
	DbName = "dms"
	DmsDbUser = "dms"
	DmsDbPassword = "d_m3123445"
}
