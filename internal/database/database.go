// Package database 提供数据库连接构建与初始化功能，支持 SQLite3 和 MySQL。
package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/babyname/fate/config"
	"github.com/babyname/fate/ent"
	"github.com/babyname/fate/ent/schema"
)

const (
	mysqlDSN   = "%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=true"
	sqlite3DSN = "file:%v?cache=shared&_journal=WAL&_fk=1"
)

// BuildFunc 根据数据库配置构建 ent 客户端的函数类型。
type BuildFunc func(config.DBConfig) (*ent.Client, error)

type database struct {
	config.DBConfig
}

// Builder 数据库构建器接口，提供获取客户端的能力。
type Builder interface {
	// Client 创建并返回数据库客户端。
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
	first, err := c.Version.Query().First(ctx)
	if err != nil && !ent.IsNotFound(err) {
		return nil, err
	}
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

// New 根据数据库配置创建 Builder 实例。
func New(cfg config.DBConfig) Builder {
	return &database{DBConfig: cfg}
}

var _ Builder = (*database)(nil)
