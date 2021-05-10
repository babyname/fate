module github.com/godcong/fate

require (
	github.com/free-utils-go/xorm_type_assist v0.0.0-20210507214645-5c65059bdc6a
	github.com/go-sql-driver/mysql v1.5.0
	github.com/godcong/chronos v0.0.3
	github.com/godcong/name_gua v0.0.0-00010101000000-000000000000
	github.com/godcong/name_wuge v0.0.0-00010101000000-000000000000
	github.com/godcong/yi v0.0.0-00010101000000-000000000000
	github.com/goextension/log v0.0.2
	github.com/mattn/go-sqlite3 v1.14.0
	github.com/spf13/cobra v1.1.3
	go.uber.org/zap v1.13.0
	xorm.io/builder v0.3.9
	xorm.io/xorm v1.0.7
)

go 1.16

// pseudo-version can be got from `go get  github.com/<name>/<project>@<commit>`
// When merged into upstream, the replace sentence can be disabled

replace (
	github.com/godcong/name_gua => github.com/fortune-fun/name_gua v0.0.0-20210510140743-082af5cba3cc
	github.com/godcong/name_wuge => github.com/fortune-fun/name_wuge v0.0.0-20210510141111-8cee898249c6
	github.com/godcong/yi => github.com/fortune-fun/yi v0.0.0-20210510181954-0bb1cb5ea507
)
