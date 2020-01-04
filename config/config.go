package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

const DefaultJSONName = "config.json"

type Config struct {
	HardMode     bool
	StrokeMax    int
	StrokeMin    int
	FixBazi      bool     //八字修正
	SupplyFilter bool     //过滤补八字
	ZodiacFilter bool     //过滤生肖
	BaguaFilter  bool     //过滤卦象
	Database     Database `json:"database"`
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

func setConfig(config *Config, src *Config) *Config {
	if config == nil {
		config = DefaultConfig()
	}
	//TODO:
	return config
}

func LoadConfig() (c *Config) {
	c = &Config{}
	def := DefaultConfig()
	f := filepath.Join(DefaultJSONPath, DefaultJSONName)
	bys, e := ioutil.ReadFile(f)
	if e != nil {
		return def
	}
	e = json.Unmarshal(bys, &c)
	if e != nil {
		return def
	}
	setConfig(c, def)
	return c
}

func DefaultConfig() *Config {
	return &Config{
		HardMode:     false,
		StrokeMax:    0,
		StrokeMin:    0,
		FixBazi:      false,
		SupplyFilter: false,
		ZodiacFilter: false,
		BaguaFilter:  false,
		Database: Database{
			Host:         "127.0.0.1",
			Port:         "3306",
			User:         "root",
			Pwd:          "111111",
			Name:         "fate",
			MaxIdleCon:   0,
			MaxOpenCon:   0,
			Driver:       "mysql",
			File:         "",
			Dsn:          "",
			ShowSQL:      true,
			ShowExecTime: true,
		},
	}
}
