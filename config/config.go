package config

import (
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

//dms系统db配置
var (
	DmsDbUrl      string
	DbName        string
	DmsDbUser     string
	DmsDbPassword string
)

/**
对于配置文件的初始化，本系统（DMS）及后续要操作的所有数据库公用一套初始化方法。
默认情况下，本业务系统默认库名：dms，默认配置为配置文件中的db0.
**/

type DbString struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Dbname   string //库名
}

//使用viper初始话数据库配置文件
func InitConfig() {

	viper.SetConfigFile("./config/dbconfig.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Panicln("读取config配置文件出错", err)
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		log.Println("检测到配置文件变更：", in.Name)
	})

}

//解组，根据输入的dbNum
func (dbstr *DbString) UnmarshalDbString(dbNum string) (err error) {
	if err = viper.UnmarshalKey(dbNum, &dbstr); err != nil {
		log.Panicln("config到struct解组失败", err)
	}
	return
}
