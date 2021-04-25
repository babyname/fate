module fate

require (
	github.com/go-sql-driver/mysql v1.5.0
	github.com/godcong/chronos v0.0.3
	github.com/godcong/fate v0.0.0-00010101000000-000000000000
	github.com/godcong/yi v0.0.0-00010101000000-000000000000
	github.com/goextension/log v0.0.2
	github.com/mattn/go-sqlite3 v1.14.0
	github.com/spf13/cobra v1.1.3
	go.uber.org/zap v1.13.0
	xorm.io/builder v0.3.9
	xorm.io/xorm v1.0.7
)

go 1.16

replace (
	github.com/godcong/fate => ../fate
	github.com/godcong/yi => ../yi
)
