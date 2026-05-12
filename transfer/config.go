package transfer

import (
	"encoding/json"
	"os"

	"github.com/babyname/fate/config"
)

type DatabaseConfig struct {
	SourceRaw json.RawMessage `json:"source"`
	Source    config.Config   `json:"-"`
	TargetRaw json.RawMessage `json:"target"`
	Target    config.Config   `json:"-"`
	Tables    []string        `json:"tables"`
	Limit     int             `json:"max"`
}

func ReadTransferConfig(p string) (*DatabaseConfig, error) {
	bytes, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}
	var db DatabaseConfig
	err = json.Unmarshal(bytes, &db)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(db.SourceRaw, &db.Source)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(db.TargetRaw, &db.Target)
	if err != nil {
		return nil, err
	}
	return &db, nil
}

func WriteTransferConfig(p string, db *DatabaseConfig) error {
	marshal, err := json.MarshalIndent(db, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p, marshal, 0644)
}
