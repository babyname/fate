package database

import (
	"errors"
	"fmt"
	"time"

	"github.com/babyname/fate/config"
	"github.com/babyname/fate/ent/schema"
	"golang.org/x/net/context"

	"github.com/babyname/fate/ent"
)

const (
	mysqlDSN   = "%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=true"
	sqlite3DSN = "file:%v?cache=shared&_journal=WAL&_fk=1"
)

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
	"other":   buildOther,
}

func buildOther(cfg config.DBConfig) (*ent.Client, error) {
	if cfg.DSN != "" {
		return nil, errors.New("dsn configuration must with a non-empty string")
	}
	return ent.Open(cfg.Driver, cfg.DSN)
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
	return ent.Open(cfg.Driver, link)
}

func (d *database) Client() (*ent.Client, error) {
	fn, ok := driverDSN[d.Driver]
	if !ok {
		fn, _ = driverDSN["other"]
	}
	c, err := fn(d.DBConfig)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	var cancel func()
	if d.DBConfig.Timeout != 0 {
		ctx, cancel = context.WithTimeout(context.Background(), time.Second*time.Duration(d.DBConfig.Timeout))
		defer cancel()
	}
	first, err := c.Version.Query().First(ctx)
	if err != nil && !ent.IsNotFound(err) {
		return nil, err
	}
	//fix empty version
	if first == nil {
		_, err := c.Version.Create().
			SetCurrentVersion(schema.CurrentDataVersion).
			SetUpdatedUnix(int(time.Now().Unix())).
			Save(ctx)
		if err != nil {
			return nil, err
		}
		return c, nil
	}
	if first.CurrentVersion != schema.CurrentDataVersion {
		return nil, fmt.Errorf("database version %d is not current,please get the correct version database", first.CurrentVersion)
	}
	return c, nil
}

// New creates a new database builder
// @param config.DBConfig
// @return Builder
func New(cfg config.DBConfig) Builder {
	return &database{DBConfig: cfg}
}

// implements check
var _ Builder = (*database)(nil)
