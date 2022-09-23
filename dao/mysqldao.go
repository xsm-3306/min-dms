package dao

import (
	"database/sql"
	"fmt"
	"log"
	"min-dms/config"

	_ "github.com/go-sql-driver/mysql"
)

//定义Mysql struct，作为mysqldao的recevier
type Database struct {
	Db_con         *sql.DB //mysql连接
	DatabaseSource string  //数据库连接串
}

//初始化数据库连接串
func (db *Database) InitDbSource() {

	db.DatabaseSource = fmt.Sprintf("%s:%s@%s/%s", config.DmsDbUser, config.DmsDbPassword, config.DmsDbUrl, config.DbName)
	log.Println(db.DatabaseSource)
}

//初始化mysql数据库连接
func (db *Database) OpenDb() (err error) {
	if db.Db_con != nil {
		return nil
	}
	db.Db_con, err = sql.Open("mysql", db.DatabaseSource)

	return
}

//close db连接
func (db *Database) CloseDb() {
	if db.Db_con != nil {
		db.Db_con.Close()
	}
	db.Db_con = nil
}

//Mysql struct实现一个SqlQuery功能，返回结果存在数组里，每个元素是map[string]string
func (db *Database) GetRows(sqlstr string, vals ...interface{}) (result []map[string]string, err error) {
	db.OpenDb()
	defer db.CloseDb()
	var rows *sql.Rows

	stmt, err := db.Db_con.Prepare(sqlstr)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err = stmt.Query(vals...)
	//rows, err = db.Db_con.Query(sqlstr, vals...)
	//查询结果处理
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	cols, _ := rows.Columns()
	l := len(cols)
	rawResult := make([][]byte, l)
	rowResult := make(map[string]string)
	//定义一个临时interface{}的数组，用来存rawresult的指针地址
	dest := make([]interface{}, l)
	for i := range rawResult {
		dest[i] = &rawResult[i]
	}

	for rows.Next() {
		err = rows.Scan(dest...) //scan rows到dest里，即指向rawResult
		if err == nil {
			for i, raw := range rawResult {
				key := cols[i]
				if raw == nil { //null空字段处理
					rowResult[key] = "null_val"
				} else {
					rowResult[key] = string(raw)
				}
				result = append(result, rowResult)
			}
		}
	}

	return
}
