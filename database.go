package fate

import (
	"fmt"
	"github.com/go-sql-driver/mysql"

	"github.com/xormsharp/xorm"
	"net/url"
)

const mysqlSource = "%s:%s@tcp(%s)/%s?loc=%s&charset=utf8mb4&parseTime=true"
const createDatabase = "CREATE DATABASE `%s` CHARACTER SET 'utf8mb4' COLLATE 'utf8mb4_bin'"

var _ = mysql.Config{}
var DefaultTableName = "fate"

func InitMysql(addr, name, pass string) *xorm.Engine {

	dbURL := fmt.Sprintf(mysqlSource, name, pass, addr, "", url.QueryEscape("Asia/Shanghai"))
	dbEngine, e := xorm.NewEngine("mysql", dbURL)
	if e != nil {
		log.Panicw("connect database failed", "error", e)
	}
	defer dbEngine.Close()
	sql := fmt.Sprintf(createDatabase, DefaultTableName)

	_, e = dbEngine.DB().Exec(sql)
	if e == nil {
		log.Infow("create database failed", "database", DefaultTableName)
	}
	u := fmt.Sprintf(mysqlSource, name, pass, addr, DefaultTableName, url.QueryEscape("Asia/Shanghai"))
	eng, e := xorm.NewEngine("mysql", u)
	if e != nil {
		log.Panicw("connect table failed", "error", e)
	}
	eng.ShowSQL(true)
	eng.ShowExecTime(true)
	return eng
}
