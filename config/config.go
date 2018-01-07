package config

import (
	"flag"

	"github.com/pelletier/go-toml"
)

type Config struct {
	tree *toml.Tree
}

var config *Config
var flags map[string]interface{}
var l = flag.String("l", "log.txt", "set log file path")
var c = flag.String("c", "config.toml", "set toml config file path")
var d = flag.Bool("d", false, "check debug output")

func init() {
	//loadConfig()
	flag.Parse()
	flags = make(map[string]interface{})
	flags["l"] = *l
	flags["c"] = *c
	flags["d"] = *d
}

func GetFlag(name string) interface{} {
	if v, b := flags[name]; b {
		return v
	}
	return nil
}

func GetFlagString(name string) string {
	if v, b := GetFlag(name).(string); b {
		return v
	}
	return ""
}

func GetFlagBool(name string) bool {
	if v, b := GetFlag(name).(bool); b {
		return v
	}
	return false
}

func loadConfig() {
	config = new(Config)
	config.tree, _ = toml.LoadFile(GetFlagString("c"))
}

func DefaultConfig() *Config {
	if config == nil {
		loadConfig()
	}
	return config
}

func (c *Config) GetSub(name string) Config {
	var cfg Config
	if v := c.Get(name); v != nil {
		cfg.tree, _ = v.(*toml.Tree)
	}
	return cfg
}

func (c *Config) Get(name string) interface{} {
	if c != nil && c.tree != nil {
		return c.tree.Get(name)
	}
	return Config{}
}

func (c *Config) GetString(name string) string {
	v, b := c.Get(name).(string)
	if b {
		return v
	}
	return ""
}

func (c *Config) GetBool(name string) bool {
	v, b := c.Get(name).(bool)
	if b {
		return v
	}
	return false
}

func (c *Config) GetStringWithDefault(name string, def string) string {
	v, b := c.Get(name).(string)
	if b {
		return v
	}
	return def
}
