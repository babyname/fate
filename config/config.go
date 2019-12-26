package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

const DefaultJSONName = "config.json"

type Config struct {
	Database Database `json:"database"`
}

var DefaultJSONPath = ""

func init() {
	if DefaultJSONPath == "" {
		dir, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		s, err := filepath.Abs(dir)
		if err != nil {
			panic(err)
		}
		DefaultJSONPath = s
	}
}

func setDefault(config *Config) *Config {
	if config == nil {
		config = &Config{}
	}
	//TODO:
	return config
}

func LoadConfig() (c *Config) {
	c = &Config{}
	f := filepath.Join(DefaultJSONPath, DefaultJSONName)
	bys, e := ioutil.ReadFile(f)
	if e != nil {
		return setDefault(c)
	}
	e = json.Unmarshal(bys, &c)
	if e != nil {
		return setDefault(c)
	}
	return c
}
