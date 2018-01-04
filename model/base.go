package model

import (
	"fmt"
	"net/url"
	"time"

	//_ "github.com/go-sql-driver/mysql"
	//"github.com/jinzhu/gorm"

	"github.com/go-xorm/xorm"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/satori/go.uuid"
)

type Base struct {
	ID        uuid.UUID `gorm:"primary_key;type:varchar(36)"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

var (
	db      *xorm.Engine
	mgTable []interface{}
)

func init() {
	ConnectDB()
}

func (b *Base) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("ID", uuid.NewV1())
	return nil
}

func ConnectDB() {

}

func NewMySql() (err error) {
	db, err = xorm.NewEngine("mysql", connectMysql())
	return err
}

func NewSqlite3() (err error) {
	db, err = xorm.NewEngine("sqlite3", "./fate.db")
	return err
}

func connectMysql() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s%sloc=%s&charset=utf8&parseTime=true",
		"root", "", "localhost", "3306", "fate", "?", url.QueryEscape("Asia/Shanghai"))
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
