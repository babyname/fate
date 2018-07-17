package config

import (
	"github.com/pelletier/go-toml"
	"os"
)

type Config struct {
	*toml.Tree
}

const DefaultFileName = "config.toml"

var config *Config

//var flags map[string]interface{}
//var l = flag.String("l", "log.txt", "set log file path")
//var c = flag.String("c", "config.toml", "set toml config file path")
//var d = flag.Bool("d", false, "check debug output")

func init() {
	config = DefaultConfig()
}

//func GetFlag(name string) interface{} {
//	if v, b := flags[name]; b {
//		return v
//	}
//	return nil
//}
//
//func GetFlagString(name string) string {
//	if v, b := GetFlag(name).(string); b {
//		return v
//	}
//	return ""
//}
//
//func GetFlagBool(name string) bool {
//	if v, b := GetFlag(name).(bool); b {
//		return v
//	}
//	return false
//}

//func loadConfig() {
//	config = new(Config)
//	config.tree, _ = toml.LoadFile(GetFlagString("c"))
//}

func NewConfig(name string) *Config {
	file, err := os.OpenFile(name, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil
	}
	tree, err := toml.LoadReader(file)
	if err != nil {
		return nil
	}
	return &Config{
		Tree: tree,
	}
}

func DefaultConfig() *Config {
	return NewConfig(DefaultFileName)
}

func (c *Config) GetTree(name string) *toml.Tree {
	if v, b := c.Get(name).(*toml.Tree); b {
		return v
	}
	return nil
}

func (c *Config) GetSub(name string) *Config {
	if v := c.GetTree(name); v != nil {
		return &Config{
			Tree: v,
		}
	}
	return nil
}

func (c *Config) GetString(name string) string {
	if v := c.Get(name); v != nil {
		if v1, b := v.(string); b {
			return v1
		}
	}
	return ""
}

func (c *Config) GetStringD(name string, def string) string {
	if v := c.Get(name); v != nil {
		if v1, b := v.(string); b {
			return v1
		}
	}
	return def
}

func (c *Config) GetBool(name string) bool {
	if v := c.Get(name); v != nil {
		if v1, b := v.(bool); b {
			return v1
		}
	}
	return false
}

func (c *Config) GetBoolD(name string, def bool) bool {
	if v := c.Get(name); v != nil {
		if v1, b := v.(bool); b {
			return v1
		}
	}
	return def
}
