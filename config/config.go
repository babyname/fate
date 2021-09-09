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
	OutputModeMarkdown
	OutputModeMax
)

type OutputFormat struct {
	OutputMode OutputMode `json:"output_mode"`
	Path       string     `json:"path"`
	Heads      []string   `json:"heads"`
}

type Config struct {
	FilterMode   FilterMode   `json:"filter_mode"`
	StrokeMax    int          `json:"stroke_max"`
	StrokeMin    int          `json:"stroke_min"`
	FixBazi      bool         `json:"fix_bazi"`      //八字修正
	SupplyFilter bool         `json:"supply_filter"` //过滤补八字
	ZodiacFilter bool         `json:"zodiac_filter"` //过滤生肖
	BaguaFilter  bool         `json:"bagua_filter"`  //过滤卦象
	DayanFilter  bool         `json:"dayan_filter"`  //过滤大衍
	Regular      bool         `json:"regular"`       //常用字过滤
	Database     DBConfig     `json:"database"`
	FileOutput   OutputFormat `json:"file_output"`
	Debug        bool
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
		Debug:        false,
		FilterMode:   0,
		StrokeMax:    18,
		StrokeMin:    3,
		DayanFilter:  false,
		FixBazi:      false,
		SupplyFilter: true,
		ZodiacFilter: true,
		BaguaFilter:  true,
		Regular:      true,
		Database: DBConfig{
			Driver: "mysql",
			DSN:    "root:111111@tcp(192.168.2.201:3306)/fate?charset=utf8\\u0026parseTime=true",
			Log:    "",
		},
		FileOutput: OutputFormat{
			Heads:      DefaultHeads,
			OutputMode: OutputModeLog,
			Path:       "name.txt",
		},
	}
}
