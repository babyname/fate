package fate

import (
	"fmt"
	"net/url"

	"github.com/go-sql-driver/mysql"
	"github.com/goextension/log"
	"github.com/mattn/go-sqlite3"

	"github.com/babyname/fate/config"

	"github.com/xormsharp/xorm"
)

const mysqlSource = "%s:%s@tcp(%s)/%s?loc=%s&charset=utf8mb4&parseTime=true"
const sqliteSource = "file:%v?cache=shared\u0026_journal=WAL\u0026_fk=1"
const createDatabase = "CREATE DATABASE `%s` CHARACTER SET 'utf8mb4' COLLATE 'utf8mb4_bin'"

var _ = mysql.MySQLDriver{}
var _ = sqlite3.SQLiteDriver{}

// DefaultTableName ...
var DefaultTableName = "fate"

// Database ...
type Database interface {
	Sync(v ...interface{}) error
	CountWuGeLucky() (n int64, e error)
	InsertOrUpdateWuGeLucky(lucky *WuGeLucky) (n int64, e error)
	GetCharacter(fn func(engine *xorm.Engine) *xorm.Session) (*Character, error)
	GetCharacters(fn func(engine *xorm.Engine) *xorm.Session) ([]*Character, error)
	FilterWuGe(last []*Character, wg chan<- *WuGeLucky) error
	Database() interface{}
}

type xormDatabase struct {
	*xorm.Engine
}

// Database ...
func (db *xormDatabase) Database() interface{} {
	return db.Engine
}

// FilterWuGe ...
func (db *xormDatabase) FilterWuGe(last []*Character, wg chan<- *WuGeLucky) error {
	return filterWuGe(db.Engine, last, wg)
}

// GetCharacters ...
func (db *xormDatabase) GetCharacters(fn func(engine *xorm.Engine) *xorm.Session) ([]*Character, error) {
	return getCharacters(db.Engine, fn)
}

// Sync ...
func (db *xormDatabase) Sync(v ...interface{}) error {
	return db.Engine.Sync2(v...)
}

// GetCharacter ...
func (db *xormDatabase) GetCharacter(fn func(engine *xorm.Engine) *xorm.Session) (*Character, error) {
	return getCharacter(db.Engine, fn)
}

// InsertOrUpdateWuGeLucky ...
func (db *xormDatabase) InsertOrUpdateWuGeLucky(lucky *WuGeLucky) (n int64, e error) {
	return insertOrUpdateWuGeLucky(db.Engine, lucky)
}

// CountWuGeLucky ...
func (db *xormDatabase) CountWuGeLucky() (n int64, e error) {
	return countWuGeLucky(db.Engine)
}

// InitDatabaseWithConfig ...
func InitDatabaseWithConfig(cfg config.Config) Database {
	return initDatabaseWithConfig(cfg.Database)
}

func initDatabaseWithConfig(db config.Database) Database {
	engine := initSQL(db)
	return &xormDatabase{
		Engine: engine,
	}
}

func databaseLink(database config.Database, schema bool) string {
	name := database.Name
	if !schema && database.Name != "" {
		name = database.Name
	}
	switch database.Driver {
	case "mysql":
		return fmt.Sprintf(mysqlSource, database.User, database.Pwd, database.Addr(), name, url.QueryEscape("Asia/Shanghai"))
	case "sqlite3":
		return fmt.Sprintf(sqliteSource, name)
	default:
		panic("unsupported database")
	}
}

func initSQL(database config.Database) *xorm.Engine {
	link := databaseLink(database, true)
	dbEngine, e := xorm.NewEngine(database.Driver, link)
	if e != nil {
		log.Panicw("connect database failed", "error", e)
	}
	sql := fmt.Sprintf(createDatabase, database.Name)
	_, e = dbEngine.DB().Exec(sql)
	if e == nil {
		log.Infow("create database failed", "database", DefaultTableName)
		dbEngine.Close()
		return nil
	}
	dbEngine.Close()
	link = databaseLink(database, false)
	eng, e := xorm.NewEngine(database.Driver, link)
	if e != nil {
		log.Panicw("connect table failed", "error", e)
	}
	if database.ShowSQL {
		eng.ShowSQL(true)
	}

	//if database.ShowExecTime {
	//	eng.ShowExecTime(true)
	//}
	return eng
}
