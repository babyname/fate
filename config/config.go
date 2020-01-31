package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

const JSONName = "config.json"

type FilterMode int

const (
	FilterModeNormal FilterMode = iota
	FilterModeHard
	FilterModeCustom
)

type OutputMode int

const (
	OutputModeLog OutputMode = iota
	OutputModeCSV
	OutputModelJSON
)

type FileOutput struct {
	OutputMode OutputMode
	Path       string
}

type Config struct {
	FilterMode   FilterMode `json:"filter_mode"`
	StrokeMax    int
	StrokeMin    int
	HardFilter   bool
	FixBazi      bool     //八字修正
	SupplyFilter bool     //过滤补八字
	ZodiacFilter bool     //过滤生肖
	BaguaFilter  bool     //过滤卦象
	Database     Database `json:"database"`
	FileOutput   FileOutput
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

func LoadConfig() (c *Config) {
	c = &Config{}
	def := DefaultConfig()
	f := filepath.Join(DefaultJSONPath, JSONName)
	bys, e := ioutil.ReadFile(f)
	if e != nil {
		return def
	}
	e = json.Unmarshal(bys, &c)
	if e != nil {
		return def
	}
	return c
}

func OutputConfig(config *Config) error {
	bys, e := json.Marshal(config)
	if e != nil {
		return e
	}

	return ioutil.WriteFile(filepath.Join(DefaultJSONPath, JSONName), bys, 0755)
}

func DefaultConfig() *Config {
	return &Config{
		FilterMode:   0,
		StrokeMax:    3,
		StrokeMin:    18,
		HardFilter:   false,
		FixBazi:      false,
		SupplyFilter: true,
		ZodiacFilter: true,
		BaguaFilter:  true,
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
		FileOutput: FileOutput{
			OutputMode: OutputModeLog,
			Path:       "name.txt",
		},
	}
}
