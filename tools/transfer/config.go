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
	err = db.Source.DecodeBytes(db.SourceRaw)
	if err != nil {
		return nil, err
	}
	err = db.Target.DecodeBytes(db.TargetRaw)
	if err != nil {
		return nil, err
	}
	return &db, nil
}

func WriteTransferConfig(p string, db *DatabaseConfig) error {
	marshal, err := jsongs.MarshalIndent(db, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(p, marshal, 0644)
}
