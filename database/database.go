package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var DB_Con *sql.DB

func InitDb() {
	db_str := "em:emp_emp123445@tcp(192.168.19.39:3306)/blog"
	Db_Con, err := sql.Open("mysql", db_str)
	if err != nil {
		fmt.Println("db open failed")
		return
	}
	err = Db_Con.Ping()
	if err != nil {
		fmt.Println("db ping failed")
		return
	}
	Db_Con.SetConnMaxLifetime(time.Minute * 10)
	Db_Con.SetMaxIdleConns(4)
	Db_Con.SetMaxOpenConns(8)
}
