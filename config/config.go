package config

import (
	"os"

	"github.com/tikafog/jsongs"
)

type Config interface {
	Database() Database
	Save(path string) error
}

type config struct {
	database database `json:"database" json-getter:"GetDatabase"`
	Debug    bool     `json:"debug"`
}

func (c *config) Save(path string) error {
	return saveConfig(path, c)
}

func (c *config) SetDatabase(database database) {
	c.database = database
}

func (c *config) Database() Database {
	return &c.database
}

func (c *config) GetDatabase() database {
	return c.database
}

func LoadConfig(path string) (c Config) {
	c = &config{}
	d := DefaultConfig()

	bys, e := os.ReadFile(path)
	if e != nil {
		return d
	}
	e = jsongs.Unmarshal(bys, &c)
	if e != nil {
		return d
	}
	return c
}

func saveConfig(path string, config *config) error {
	bytes, e := jsongs.MarshalIndent(config, "", " ")
	if e != nil {
		return e
	}
	return os.WriteFile(path, bytes, 0644)
}

func DefaultConfig() Config {
	return &config{
		Debug: false,
		database: database{
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

func DefaultSqliteConfig() Config {
	return &config{
		Debug: false,
		database: database{
			driver: "sqlite3",
			dsn:    sqlite3DSN,
			dbName: "fate",
		},
	}
}
