package fate

import (
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/mattn/go-sqlite3"
	"net/url"
)

const sqlite3Source = "file:%s?cache=shared&mode=rwc&_journal_mode=WAL"
const mysqlSource = "%s:%s@tcp(%s)/%s?loc=%s&charset=utf8mb4&parseTime=true"
const createDatabase = "CREATE DATABASE `%s` CHARACTER SET 'utf8mb4' COLLATE 'utf8mb4_bin'"

var _ = mysql.Config{}
var _ = sqlite3.SQLiteDriver{}
var DefaultTableName = "fate"

func NewSQLite3(name string) (eng *xorm.Engine, e error) {
	eng, e = xorm.NewEngine("sqlite3", fmt.Sprintf(sqlite3Source, name))
	if e != nil {
		return nil, e
	}
	return eng, nil
}

func LoadCharacter(path string) (eng *xorm.Engine, e error) {
	eng, e = xorm.NewEngine("sqlite3", fmt.Sprintf(sqlite3Source, path))
	if e != nil {
		return nil, e
	}
	return eng, nil
}

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
