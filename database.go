package fate

import (
	"fmt"
	"net/url"
	_ "time/tzdata"

	"github.com/go-sql-driver/mysql"
	"github.com/godcong/fate/config"
	"github.com/goextension/log"
	"github.com/mattn/go-sqlite3"
	"xorm.io/xorm"
)

var _ = &sqlite3.SQLiteDriver{}

const mysqlSource = "%s:%s@tcp(%s)/%s?loc=%s&charset=utf8mb4&parseTime=true"
const createMysqlDatabase = "CREATE DATABASE `%s` CHARACTER SET 'utf8mb4' COLLATE 'utf8mb4_unicode_ci'"

var _ = mysql.Config{}

// DefaultTableName ...
var DefaultTableName = "fate"

// 持有*xorm.Engine并增加方法
type Database interface {
	//重命名xorm.Engine.Sync2方法
	Sync(v ...interface{}) error
	//用查询语句查出一个字符
	GetCharacter(fn func(engine *xorm.Engine) *xorm.Session) (*Character, error)
	//用查询语句查出符合条件的所有字符
	GetCharacters(fn func(engine *xorm.Engine) *xorm.Session) ([]*Character, error)
	//返回持有的*xorm.Engine
	Database() interface{}
}

//数据库方法实现类
type xormDatabase struct {
	*xorm.Engine
}

// Database ...
func (db *xormDatabase) Database() interface{} {
	return db.Engine
}

// GetCharacters ...
func (db *xormDatabase) GetCharacters(fn func(engine *xorm.Engine) *xorm.Session) ([]*Character, error) {
	return getCharacters(db.Engine, fn)
}

//用查询语句查出符合条件的所有字符
func getCharacters(engine *xorm.Engine, fn func(engine *xorm.Engine) *xorm.Session) ([]*Character, error) {
	s := fn(engine)
	var c []*Character
	e := s.Find(&c)
	if e == nil {
		return c, nil
	}
	return nil, fmt.Errorf("character list get error:%w", e)
}

// Sync ...
func (db *xormDatabase) Sync(v ...interface{}) error {
	return db.Engine.Sync2(v...)
}

// GetCharacter ...
func (db *xormDatabase) GetCharacter(fn func(engine *xorm.Engine) *xorm.Session) (*Character, error) {
	return getCharacter(db.Engine, fn)
}

//用查询语句查出一个字符
func getCharacter(eng *xorm.Engine, fn func(engine *xorm.Engine) *xorm.Session) (*Character, error) {
	s := fn(eng)
	var c Character
	b, e := s.Get(&c)
	if e == nil && b {
		return &c, nil
	}
	return nil, fmt.Errorf("character get error:%w", e)
}

// InitDatabaseWithConfig ...
func InitDatabaseWithConfig(cfg config.Config) Database {
	var engine *xorm.Engine

	if cfg.Database.Driver == "mysql" {
		engine = initMysql(cfg.Database)
	} else if cfg.Database.Driver == "sqlite3" {
		engine = initSqlite(cfg.Database)
	} else {
		panic("engine type not supported")
	}

	return &xormDatabase{
		Engine: engine,
	}
}

func initSqlite(dbCfg config.DatabaseConfig) *xorm.Engine {
	eng, e := xorm.NewEngine(dbCfg.Driver, dbCfg.File)
	if e != nil {
		log.Panicw("connect table failed", "error", e)
	}
	if dbCfg.ShowSQL {
		eng.ShowSQL(true)
	}

	_, e = eng.Exec("PRAGMA journal_mode = OFF;")
	if e != nil {
		log.Fatal(e)
	}

	return eng
}

func initMysql(database config.DatabaseConfig) *xorm.Engine {
	dbURL := fmt.Sprintf(mysqlSource, database.User, database.Pwd, database.Addr(), "", url.QueryEscape("Asia/Shanghai"))
	dbEngine, e := xorm.NewEngine(database.Driver, dbURL)
	if e != nil {
		log.Panicw("connect database failed", "error", e)
	}
	defer dbEngine.Close()

	sql := fmt.Sprintf(createMysqlDatabase, database.Name)

	_, e = dbEngine.DB().Exec(sql)
	if e == nil {
		log.Infow("create database failed", "database", DefaultTableName)
	}
	u := fmt.Sprintf(mysqlSource, database.User, database.Pwd, database.Addr(), database.Name, url.QueryEscape("Asia/Shanghai"))
	eng, e := xorm.NewEngine(database.Driver, u)
	if e != nil {
		log.Panicw("connect table failed", "error", e)
	}
	if database.ShowSQL {
		eng.ShowSQL(true)
	}

	return eng
}
