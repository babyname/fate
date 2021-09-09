package model

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/goextension/log"

	"github.com/godcong/fate/config"
	"github.com/godcong/fate/ent"
)

const dsn = "%s:%s@tcp(%s)/%s?loc=%s&charset=utf8mb4&parseTime=true"

type Model struct {
	*ent.Client
	cfg config.DBConfig
}

func New(cfg config.Config) (*Model, error) {
	var options []ent.Option
	if cfg.Debug {
		options = append(options, ent.Debug())
	}

	if cfg.Database.Log != "" {
		options = append(options, ent.Log(func(i ...interface{}) {
			log.Debug(i...)
		}))
	}

	open, err := ent.Open(cfg.Database.Driver, cfg.Database.DSN, options...)
	if err != nil {
		return nil, err
	}
	return &Model{
		Client: open,
		cfg:    cfg.Database,
	}, nil
}

// Hash ...
func Hash(url string) string {
	s := sha256.New()
	return hex.EncodeToString(s.Sum([]byte(url)))
}
