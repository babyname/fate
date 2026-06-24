package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/babyname/fate/v4/config"
	"github.com/babyname/fate/v4/ent"
	_ "github.com/sqlite3ent/sqlite3"
)

const (
	mysqlDSN = "%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=true"
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
	if cfg.DSN != "" {
		return ent.Open(cfg.Driver, cfg.DSN)
	}
	switch cfg.Mode {
	case "memory":
		return ent.Open(cfg.Driver, "file::memory:?cache=shared&_fk=1")
	case "file", "":
		initMode := cfg.InitMode
		if initMode == "" {
			initMode = config.InitModeAuto
		}

		switch initMode {
		case config.InitModeDB:
			dbFile, err := ensureDBFile(cfg)
			if err != nil {
				return nil, err
			}
			if dbFile == "" {
				return nil, fmt.Errorf("database file not found and init_mode is 'db'")
			}
			return ent.Open(cfg.Driver, fmt.Sprintf("file:%s?cache=shared&_journal=WAL&_fk=1", dbFile))
		case config.InitModeJSON:
			dbFile := cfg.GetDBFile()
			return ent.Open(cfg.Driver, fmt.Sprintf("file:%s?cache=shared&_journal=WAL&_fk=1", dbFile))
		case config.InitModeAuto:
			fallthrough
		default:
			dbFile, err := ensureDBFile(cfg)
			if err != nil {
				return nil, err
			}
			if dbFile != "" {
				return ent.Open(cfg.Driver, fmt.Sprintf("file:%s?cache=shared&_journal=WAL&_fk=1", dbFile))
			}
			dbFile = cfg.GetDBFile()
			return ent.Open(cfg.Driver, fmt.Sprintf("file:%s?cache=shared&_journal=WAL&_fk=1", dbFile))
		}
	default:
		return nil, fmt.Errorf("unknown sqlite3 mode: %s", cfg.Mode)
	}
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
		fn = driverDSN["other"]
	}
	c, err := fn(d.DBConfig)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	var cancel func()
	if d.Timeout != 0 {
		ctx, cancel = context.WithTimeout(context.Background(), time.Second*time.Duration(d.Timeout))
		defer cancel()
	}
	if err := c.Schema.Create(ctx); err != nil {
		return nil, fmt.Errorf("create schema: %w", err)
	}

	initMode := d.InitMode
	if initMode == "" {
		initMode = config.InitModeAuto
	}

	needs, err := needsInit(c)
	if err != nil {
		return nil, err
	}

	if !needs {
		return c, nil
	}

	switch initMode {
	case config.InitModeDB:
		if d.Driver == "sqlite3" {
			return nil, fmt.Errorf("database version mismatch and init_mode is 'db'")
		}
		if err := initializeFromJSON(ctx, c); err != nil {
			return nil, fmt.Errorf("initialize from json: %w", err)
		}
	case config.InitModeJSON:
		if err := initializeFromJSON(ctx, c); err != nil {
			return nil, fmt.Errorf("initialize from json: %w", err)
		}
	case config.InitModeAuto:
		fallthrough
	default:
		if err := initializeFromJSON(ctx, c); err != nil {
			return nil, fmt.Errorf("initialize from json: %w", err)
		}
	}

	return c, nil
}

func New(cfg config.DBConfig) Builder {
	return &database{DBConfig: cfg}
}

var _ Builder = (*database)(nil)
