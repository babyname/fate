package database

import (
	"fmt"
	"github.com/babyname/fate/config"

	"github.com/babyname/fate/ent"
)

const mysqlDSN = "%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=true"
const sqlite3DSN = "file:%v?cache=shared&_journal=WAL&_fk=1"

type BuildFunc func(config.DBConfig) (*ent.Client, error)

type database struct {
	config.DBConfig
}

type Builder interface {
	Client() (*ent.Client, error)
}

var driverDSN = map[string]BuildFunc{
	"sqlite3": buildSqlite3,
	"mysql":   buildMysql,
}

func buildSqlite3(cfg config.DBConfig) (*ent.Client, error) {
	dsn := sqlite3DSN
	if cfg.DSN != "" {
		dsn = cfg.DSN
	}
	link := fmt.Sprintf(dsn, cfg.Name)
	return ent.Open(cfg.Driver, link)
}

func buildMysql(cfg config.DBConfig) (*ent.Client, error) {
	dsn := mysqlDSN
	if cfg.DSN != "" {
		dsn = cfg.DSN
	}
	link := fmt.Sprintf(dsn, cfg.User, cfg.Pwd, cfg.Host, cfg.Port, cfg.Name)
	fmt.Println("open:", link)
	return ent.Open(cfg.Driver, link)
}

func (d *database) Client() (*ent.Client, error) {
	fn, ok := driverDSN[d.Driver]
	if !ok {
		return nil, fmt.Errorf("the driver of %v is not supported", d.Driver)
	}
	return fn(d.DBConfig)
}

// New creates a new database builder
// @param config.DBConfig
// @return Builder
func New(cfg config.DBConfig) Builder {
	return &database{DBConfig: cfg}
}

// implements check
var _ Builder = (*database)(nil)
