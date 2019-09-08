package fate

import (
	"fmt"
	"github.com/go-xorm/xorm"
)

const SQLite3Source = "file:%s?cache=shared&mode=rwc&_journal_mode=WAL"

func NewSQLite3(name string) (eng *xorm.Engine, e error) {
	eng, e = xorm.NewEngine("sqlite3", fmt.Sprintf(SQLite3Source, name))
	if e != nil {
		return nil, e
	}
	return eng, nil
}

type luckyFilter struct {
	fate  *fate
	limit int
	start int
}

func AllLucy(session *xorm.Session, l1, l2 int) ([]*WuGeLucky, error) {
	llist := new([]*WuGeLucky)
	e := session.Where("last_stroke_1 = ?", l1).And("last_stroke_2 = ?", l2).Find(llist)
	if e != nil {
		return nil, e
	}
	return *llist, nil
}
