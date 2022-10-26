package dao

import (
	"database/sql"
	"fmt"
	"min-dms/config"

	_ "github.com/go-sql-driver/mysql"
)

//定义Mysql struct，作为mysqldao的recevier
type Database struct {
	Db_con         *sql.DB //mysql连接
	DatabaseSource string  //数据库连接串
}

//初始化数据库连接串
func (db *Database) InitDbSource(dbNum string, dbName string) {
	var dbstr config.DbString
	dbstr.UnmarshalDbString(dbNum)
	//log.Println(dbstr)
	dbstr.Dbname = dbName
	db.DatabaseSource = fmt.Sprintf("%s:%s@tcp(%s:%v)/%s", dbstr.Username, dbstr.Password, dbstr.Host, dbstr.Port, dbstr.Dbname)
	//log.Println("initdbsource:", db.DatabaseSource)
}

func (db *Database) NewDb(dbNum string, dbName string) *Database {
	var newdb = new(Database)
	newdb.InitDbSource(dbNum, dbName)
	return newdb
}

//初始化mysql数据库连接
func (db *Database) OpenDb() error {

	if db.Db_con != nil {
		return nil
	}
	var err error
	db.Db_con, err = sql.Open("mysql", db.DatabaseSource)

	return err
}

//close db连接
func (db *Database) CloseDb() {
	if db.Db_con != nil {
		db.Db_con.Close()
	}
	db.Db_con = nil
}
