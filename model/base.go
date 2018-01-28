package model

import (
	"fmt"
	"net/url"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/godcong/fate/config"
	"github.com/godcong/fate/debug"
	_ "github.com/mattn/go-sqlite3"
	uuid "github.com/satori/go.uuid"
)

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
		db.GetStringWithDefault("username", "root"),
		db.GetStringWithDefault("password", ""),
		db.GetStringWithDefault("addr", "localhost"),
		db.GetStringWithDefault("port", "3306"),
		db.GetStringWithDefault("schema", "default"),
		db.GetStringWithDefault("param", "?"),
		url.QueryEscape(db.GetStringWithDefault("local", "Asia/Shanghai")))
}

func ConnectDB(config *config.Config) *xorm.Engine {
	database := config.GetSub("database")
	driver := database.GetStringWithDefault("name", "sqlite3")
	source := ""
	if driver == "mysql" {
		source = connectMysql(config)
	} else if driver == "sqlite3" {
		source = database.GetStringWithDefault("path", "fate")
	}
	if NewDatabase(driver, source) != nil {
		return nil
	}
	return db
}

func NewDatabase(driver, source string) (err error) {
	db, err = xorm.NewEngine(driver, source)
	if config.DefaultConfig().GetBool("system.sql") == true {
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

//func CreateDB() *DB {
//	db, e := gorm.Open("mysql", connectString())
//	if e != nil {
//		panic(e)
//	}
//	return db
//
//}

//func ORM() *DB {
//	if db == nil || db.DB().Ping() != nil {
//		db = CreateDB()
//	}
//	db.LogMode(true)
//	return db
//}

//func SetMigrate(table interface{}) {
//	migrate = append(migrate, table)
//
//}

//func RunMigrate() {
//	for _, v := range migrate {
//		ORM().AutoMigrate(v)
//	}
//}

//func (b Base) IsNil() bool {
//	return uuid.Equal(b.ID, uuid.Nil)
//}
