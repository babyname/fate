package config

import (
	"os"

	"github.com/tikafog/jsongs"
)

type Config struct {
	Database database `json:"database"`
	Debug    bool     `json:"debug"`
}

func (c Config) GetDatabase() Database {
	return &c.Database
}

func LoadConfig(path string) (c *Config) {
	c = &Config{}
	def := DefaultConfig()

	bys, e := os.ReadFile(path)
	if e != nil {
		return def
	}
	e = jsongs.Unmarshal(bys, &c)
	if e != nil {
		return def
	}
	return c
}

func SaveConfig(path string, config *Config) error {
	bytes, e := jsongs.MarshalIndent(config, "", " ")
	if e != nil {
		return e
	}
	return os.WriteFile(path, bytes, 0644)
}

func DefaultConfig() *Config {
	return &Config{
		Debug: false,
		Database: database{
			driver: "mysql",
			dsn:    mysqlDSN,
			host:   "localhost",
			port:   "3306",
			user:   "root",
			pwd:    "root",
			dbName: "fate",
		},
	}
}

func DefaultSqliteConfig() *Config {
	return &Config{
		Debug: false,
		Database: database{
			driver: "sqlite3",
			dsn:    sqlite3DSN,
			dbName: "fate",
		},
	}
}
