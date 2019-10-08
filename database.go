package fate

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"net/url"
)

const SQLite3Source = "file:%s?cache=shared&mode=rwc&_journal_mode=WAL"

func NewSQLite3(name string) (eng *xorm.Engine, e error) {
	eng, e = xorm.NewEngine("sqlite3", fmt.Sprintf(SQLite3Source, name))
	if e != nil {
		return nil, e
	}
	return eng, nil
}

func LoadCharacter(path string) (eng *xorm.Engine, e error) {
	eng, e = xorm.NewEngine("sqlite3", fmt.Sprintf(SQLite3Source, path))
	if e != nil {
		return nil, e
	}
	return eng, nil
}

const sqlURL = "%s:%s@tcp(%s)/%s?loc=%s&charset=utf8mb4&parseTime=true"

func InitMysql(addr, name, pass string) *xorm.Engine {
	u := fmt.Sprintf(sqlURL, name, pass, addr, "fate", url.QueryEscape("Asia/Shanghai"))
	eng, e := xorm.NewEngine("mysql", u)
	if e != nil {
		log.Fatal(e)
	}
	eng.ShowSQL(true)
	eng.ShowExecTime(true)
	return eng
}
