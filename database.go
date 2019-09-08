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
