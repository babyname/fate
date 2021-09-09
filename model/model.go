package model

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"

	"github.com/godcong/fate/config"
	"github.com/godcong/fate/ent"
)

const dsn = "%s:%s@tcp(%s)/%s?loc=%s&charset=utf8mb4&parseTime=true"

type Model struct {
	*ent.Client
	cfg config.DBConfig
}

func New(cfg config.DBConfig) (*Model, error) {
	dsnPath := fmt.Sprintf(dsn, cfg.User, cfg.Pwd, cfg.Addr, cfg.Name, url.QueryEscape(cfg.Loc))

	var options []ent.Option
	if cfg.Debug {
		options = append(options, ent.Debug())
	}

	if cfg.Log != "" {
		options = append(options, ent.Log(func(i ...interface{}) {
			//todo:
		}))
	}

	open, err := ent.Open(cfg.Driver, dsnPath, options...)
	if err != nil {
		return nil, err
	}
	return &Model{
		Client: open,
		cfg:    cfg,
	}, nil
}

// Hash ...
func Hash(url string) string {
	s := sha256.New()
	return hex.EncodeToString(s.Sum([]byte(url)))
}
