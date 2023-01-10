package config

import (
	"encoding/json"
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

	//RunInit      bool       `json:"run_init"`
	//FilterMode   FilterMode `json:"filter_mode"`
	//StrokeMax    int        `json:"stroke_max"`
	//StrokeMin    int        `json:"stroke_min"`
	//HardFilter   bool       `json:"hard_filter"`
	//FixBazi      bool       `json:"fix_bazi"`      //八字修正
	//SupplyFilter bool       `json:"supply_filter"` //过滤补八字
	//ZodiacFilter bool       `json:"zodiac_filter"` //过滤生肖
	//BaguaFilter  bool       `json:"bagua_filter"`  //过滤卦象
	//Regular      bool       `json:"regular"`       //常用
	//FileOutput   FileOutput `json:"file_output"`
	WorkingDir string    `json:"working_dir,omitempty"`
	Database   DBConfig  `json:"database"`
	Log        LogConfig `json:"logger,omitempty"`
}

const (
	defaultConfigPath = ""
	defaultWorkingDir = "fate"
)

var DefaultHeads = []string{"姓名", "笔画", "拼音", "喜用神", "八字"}

func LoadConfig(path string) (c *Config) {
	c = &Config{}
	def := DefaultConfig()
	f := filepath.Join(defaultConfigPath, JSONName)
	bys, e := os.ReadFile(f)
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

	return os.WriteFile(filepath.Join(defaultConfigPath, JSONName), bys, 0755)
}

func DefaultConfig() *Config {
	return &Config{
		WorkingDir: defaultWorkingDir,
		Database:   defaultDBSqlite3(),
		Log:        defaultLogConfig(),
	}
}

func GetPath(root string, paths ...string) string {
	workpath := filepath.Join(paths...)
	path := filepath.Join(defaultWorkingDir, workpath)
	if root == "" {
		dir, err := os.Getwd()
		if err != nil {
			return path
		}
		s, err := filepath.Abs(dir)
		if err != nil {
			return path
		}
		return s
	}
	return filepath.Join(root, workpath)
}
