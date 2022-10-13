package transfer

import (
	"encoding/json"
	"os"

	"github.com/tikafog/jsongs"

	"github.com/babyname/fate/config"
)

type DatabaseConfig struct {
	SourceRaw json.RawMessage `json:"source"`
	Source    config.Config   `json:"-"`
	TargetRaw json.RawMessage `json:"target"`
	Target    config.Config   `json:"-"`
	Tables    []string
	Limit     int `json:"max"`
}

func ReadTransferConfig(p string) (*DatabaseConfig, error) {
	bytes, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}
	var db DatabaseConfig
	err = jsongs.Unmarshal(bytes, &db)
	if err != nil {
		return nil, err
	}
	db.Source, err = config.LoadFromBytes(db.SourceRaw)
	if err != nil {
		return nil, err
	}
	db.Target, err = config.LoadFromBytes(db.TargetRaw)
	return &db, nil
}

func WriteTransferConfig(p string, db *DatabaseConfig) error {
	marshal, err := jsongs.MarshalIndent(db, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(p, marshal, 0644)
}
