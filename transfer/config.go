package transfer

import (
	"os"

	"github.com/tikafog/jsongs"

	"github.com/babyname/fate/config"
)

type DatabaseConfig struct {
	Source *config.Config `json:"source"`
	Target *config.Config `json:"target"`
	Tables []string
}

func readConfig(p string) (*DatabaseConfig, error) {
	bytes, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}
	var db DatabaseConfig
	err = jsongs.Unmarshal(bytes, &db)
	if err != nil {
		return nil, err
	}
	return &db, nil
}

func writeConfig(p string, db *DatabaseConfig) error {
	marshal, err := jsongs.MarshalIndent(db, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(p, marshal, 0644)
}
