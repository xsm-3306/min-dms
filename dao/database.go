package dao

import (
	"database/sql"
	"fmt"
	"log"
	"min-dms/config"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Db struct {
	Dbstring string
	Con_pool *sql.DB
}

func (db *Db) Init() {
	db.Dbstring = fmt.Sprintf("%s:%s@%s/%s", config.User, config.Password, config.Host, config.Db_name)
}

func (db *Db) Open() error {
	if db.Con_pool != nil {
		return nil
	}
	var err error
	db.Con_pool, err = sql.Open("mysql", db.Dbstring)
	return err
}

var Con_pool *sql.DB

//使用golang原生sql包，初始话数据库连接
//此处initdb初始话为dms系统的业务配置库
func InitDb() {
	var err error
	db_str := "dms:d_m3123445@tcp(192.168.19.39:3306)/dms"
	Con_pool, err = sql.Open("mysql", db_str)
	if err != nil {
		log.Fatal(err)
		return
	}
	err = Con_pool.Ping()
	if err != nil {
		log.Fatal(err)
		return
	}
	Con_pool.SetConnMaxLifetime(time.Minute * 10)
	Con_pool.SetMaxIdleConns(4)
	Con_pool.SetMaxOpenConns(8)
}
