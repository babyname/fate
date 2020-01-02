package fate

import (
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/godcong/fate/config"
	"github.com/goextension/log"

	"github.com/xormsharp/xorm"
	"net/url"
)

const mysqlSource = "%s:%s@tcp(%s)/%s?loc=%s&charset=utf8mb4&parseTime=true"
const createDatabase = "CREATE DATABASE `%s` CHARACTER SET 'utf8mb4' COLLATE 'utf8mb4_bin'"

var _ = mysql.Config{}
var DefaultTableName = "fate"

type Database interface {
}

type xormDatabase struct {
	*xorm.Engine
}

func InitFromConfig(config config.Config) Database {
	engine := initSQL(config.Database)
	return xormDatabase{
		Engine: engine,
	}
}

func initSQL(database config.Database) *xorm.Engine {
	dbURL := fmt.Sprintf(mysqlSource, database.User, database.Pwd, database.Addr(), "", url.QueryEscape("Asia/Shanghai"))
	dbEngine, e := xorm.NewEngine(database.Driver, dbURL)
	if e != nil {
		log.Panicw("connect database failed", "error", e)
	}
	defer dbEngine.Close()

	sql := fmt.Sprintf(createDatabase, database.Name)

	_, e = dbEngine.DB().Exec(sql)
	if e == nil {
		log.Infow("create database failed", "database", DefaultTableName)
	}
	u := fmt.Sprintf(mysqlSource, database.User, database.Pwd, database.Addr(), database.Name, url.QueryEscape("Asia/Shanghai"))
	eng, e := xorm.NewEngine("mysql", u)
	if e != nil {
		log.Panicw("connect table failed", "error", e)
	}
	if database.ShowSQL {
		eng.ShowSQL(true)
	}

	if database.ShowExecTime {
		eng.ShowExecTime(true)
	}
	return eng
}
