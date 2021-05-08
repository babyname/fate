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
	RunInit      bool           `json:"run_init"`
	FilterMode   FilterMode     `json:"filter_mode"`
	StrokeMax    int            `json:"stroke_max"`
	StrokeMin    int            `json:"stroke_min"`
	HardFilter   bool           `json:"hard_filter"`
	FixBazi      bool           `json:"fix_bazi"`      //八字修正
	SupplyFilter bool           `json:"supply_filter"` //过滤补八字
	ZodiacFilter bool           `json:"zodiac_filter"` //过滤生肖
	BaguaFilter  bool           `json:"bagua_filter"`  //过滤卦象
	Regular      bool           `json:"regular"`       //常用
	Database     DatabaseConfig `json:"database"`
	FileOutput   FileOutput     `json:"file_output"`
}

var DefaultJSONPath = ""
var DefaultHeads = []string{"姓名", "笔画", "拼音", "八字"}

//"喜用神", 等待喜用神和调和用神完善后再加入

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
	f := filepath.Join(DefaultJSONPath, JSONName)
	bys, e := ioutil.ReadFile(f)
	if e != nil {
		panic(f)
	}
	e = json.Unmarshal(bys, &c)
	if e != nil {
		panic(bys)
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
		StrokeMax:    24, //输出最大笔画数
		StrokeMin:    1,  //输出最小笔画数
		HardFilter:   false,
		FixBazi:      false, //立春修正（出生日期为立春当日时间为已过立春八字需修正）
		SupplyFilter: true,  //三才五格过滤
		ZodiacFilter: true,  //生肖过滤
		BaguaFilter:  true,  //周易八卦过滤
		Regular:      true,
		Database: DatabaseConfig{ //连接DB：
			Host:         "localhost",
			Port:         "3306",
			User:         "root",
			Pwd:          "111111",
			Name:         "fate",
			MaxIdleCon:   0,
			MaxOpenCon:   0,
			Driver:       "sqlite3",
			File:         "fate.db",
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
