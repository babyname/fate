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
	OutputMode OutputMode `json:"output_mode"`
	Path       string     `json:"path"`
	Heads      []string   `json:"heads"`
}

type Config struct {
	RunInit      bool       `json:"run_init"`
	FilterMode   FilterMode `json:"filter_mode"`
	StrokeMax    int        `json:"stroke_max"`
	StrokeMin    int        `json:"stroke_min"`
	HardFilter   bool       `json:"hard_filter"`
	FixBazi      bool       `json:"fix_bazi"`      //八字修正
	SupplyFilter bool       `json:"supply_filter"` //过滤补八字
	ZodiacFilter bool       `json:"zodiac_filter"` //过滤生肖
	BaguaFilter  bool       `json:"bagua_filter"`  //过滤卦象
	Regular      bool       `json:"regular"`       //常用
	Database     Database   `json:"database"`
	FileOutput   FileOutput `json:"file_output"`
}

var DefaultJSONPath = ""
var DefaultHeads = []string{"姓名", "笔画", "拼音", "喜用神", "八字"}

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
	bys, e := json.MarshalIndent(config, "", " ")
	if e != nil {
		return e
	}

	return ioutil.WriteFile(filepath.Join(DefaultJSONPath, JSONName), bys, 0755)
}

func DefaultConfig() *Config {
	return &Config{
		RunInit:      false,
		FilterMode:   0,
		StrokeMax:    18,
		StrokeMin:    3,
		HardFilter:   false,
		FixBazi:      false,
		SupplyFilter: true,
		ZodiacFilter: true,
		BaguaFilter:  true,
		Regular:      true,
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
			ShowSQL:      false,
			ShowExecTime: false,
		},
		FileOutput: FileOutput{
			Heads:      DefaultHeads,
			OutputMode: OutputModeLog,
			Path:       "name.txt",
		},
	}
}
