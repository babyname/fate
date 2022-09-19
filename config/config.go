package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Database Database `json:"database"`
	Debug    bool     `json:"debug"`
}

func LoadConfig(path string) (c *Config) {
	c = &Config{}
	def := DefaultConfig()

	bys, e := os.ReadFile(path)
	if e != nil {
		return def
	}
	e = json.Unmarshal(bys, &c)
	if e != nil {
		return def
	}
	return c
}

func DefaultConfig() *Config {
	return &Config{
		Debug: false,
		Database: Database{
			Driver: "mysql",
			DSN:    "root:111111@tcp(127.0.0.1:3306)/fate?charset=utf8\\u0026parseTime=true",
		},
	}
}
