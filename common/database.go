package common

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var Con_pool *sql.DB

//使用golang原生sql包，初始话数据库连接
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
