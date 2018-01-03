package model

import (
	"time"

	"fmt"

	"net/url"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
)

type Base struct {
	ID        uuid.UUID `gorm:"primary_key;type:varchar(36)"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

type DB = gorm.DB

var (
	db      *DB
	migrate = []interface{}{}
)

func init() {
	//db = CreateDB()
}

func (b *Base) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("ID", uuid.NewV1())
	return nil
}

func connectString() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s%sloc=%s&charset=utf8&parseTime=true",
		"root", "", "localhost", "3306", "fate", "?", url.QueryEscape("Asia/Shanghai"))
}

func CreateDB() *DB {
	db, e := gorm.Open("mysql", connectString())
	if e != nil {
		panic(e)
	}
	return db

}

func ORM() *DB {
	if db == nil || db.DB().Ping() != nil {
		db = CreateDB()
	}
	db.LogMode(true)
	return db
}

func SetMigrate(table interface{}) {
	migrate = append(migrate, table)

}

func RunMigrate() {
	for _, v := range migrate {
		ORM().AutoMigrate(v)
	}
}

func (b Base) IsNil() bool {
	return uuid.Equal(b.ID, uuid.Nil)
}
