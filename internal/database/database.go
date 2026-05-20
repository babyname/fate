package database

import (
	"bytes"
	"compress/gzip"
	"context"
	_ "embed"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/babyname/fate/config"
	"github.com/babyname/fate/ent"
	"github.com/babyname/fate/ent/schema"
	_ "github.com/sqlite3ent/sqlite3"
)

//go:embed data/fate.db.gz
var embeddedDB []byte

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
		name := cfg.Name
		if name == "" {
			name = "fate"
		}
		if HasEmbeddedDB() {
			if err := ExtractEmbeddedDB(name); err != nil {
				return nil, fmt.Errorf("extract embedded db: %w", err)
			}
		}
		return ent.Open(cfg.Driver, fmt.Sprintf("file:%s?cache=shared&_journal=WAL&_fk=1", name))
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

func New(cfg config.DBConfig) Builder {
	return &database{DBConfig: cfg}
}

func HasEmbeddedDB() bool {
	return len(embeddedDB) > 0
}

func ExtractEmbeddedDB(destPath string) error {
	if _, err := os.Stat(destPath); err == nil {
		return nil
	}
	log.Printf("[DB] Extracting embedded database to %s (%d bytes compressed)", destPath, len(embeddedDB))
	gr, err := gzip.NewReader(bytes.NewReader(embeddedDB))
	if err != nil {
		return fmt.Errorf("create gzip reader: %w", err)
	}
	defer gr.Close()

	out, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("create db file: %w", err)
	}
	defer out.Close()

	written, err := io.Copy(out, gr)
	if err != nil {
		os.Remove(destPath)
		return fmt.Errorf("decompress db: %w", err)
	}
	log.Printf("[DB] Extracted database: %d bytes", written)
	return nil
}

var _ Builder = (*database)(nil)
