package model

import (
	"fmt"
	"net/url"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/godcong/fate/config"
	"github.com/godcong/fate/debug"
	"github.com/satori/go.uuid"
)

type SyncAble interface {
	Sync() error
}

type Base struct {
	Id      string     `xorm:"uuid pk"`
	Created time.Time  `xorm:"created"`
	Updated time.Time  `xorm:"updated"`
	Deleted *time.Time `xorm:"deleted"`
	Version int        `xorm:"version"`
}

type Model interface {
	Create(v ...interface{}) (int64, error)
}

var (
	db      *xorm.Engine
	mgTable []interface{}
)

func init() {
	ConnectDB(config.DefaultConfig())
}

func (b *Base) BeforeInsert() {
	b.Id = uuid.NewV1().String()
}

func connectMysql(config *config.Config) string {
	db := config.GetSub("database")
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s%sloc=%s&charset=utf8&parseTime=true",
		db.GetStringD("username", "root"),
		db.GetStringD("password", "111111"),
		db.GetStringD("addr", "localhost"),
		db.GetStringD("port", "3306"),
		db.GetStringD("schema", "fate"),
		db.GetStringD("param", "?"),
		url.QueryEscape(db.GetStringD("local", "Asia/Shanghai")))
}

func ConnectDB(config *config.Config) *xorm.Engine {
	database := config.GetSub("database")
	driver := database.GetStringD("name", "mysql")
	source := ""
	if driver == "mysql" {
		source = connectMysql(config)
	} else {
		panic("no sql server")
	}
	if NewDatabase(driver, source) != nil {
		return nil
	}
	db.ShowSQL(true)
	return db
}

func NewDatabase(driver, source string) (err error) {
	db, err = xorm.NewEngine(driver, source)
	if config.DefaultConfig().GetBool("system.show_sql") == true {
		db.ShowSQL(true)
	}
	return err
}

var tables []interface{}

func Register(i interface{}) {
	tables = append(tables, i)
}

func CreateTables() error {
	debug.Println("CreateTables")
	return db.CreateTables(tables...)
}

func Sync(v interface{}) error {
	return db.Sync2(v)
}

func DB() *xorm.Engine {
	if db == nil || db.DB().Ping() != nil {
		db = ConnectDB(config.DefaultConfig())
	}
	return db
}
