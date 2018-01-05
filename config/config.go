package config

import (
	"flag"
	"log"

	"github.com/pelletier/go-toml"
)

type Config struct {
	//Database []map[string]string `toml:"database"`
	tree *toml.Tree
}

var config *Config

func init() {
	//loadConfig()
}

func loadConfig() {
	//bs, err := ioutil.ReadFile(*c)
	//if err != nil {
	//	log.Panic(err)
	//}
	var err error
	c := flag.String("c", "config.toml", "load an toml config file")
	config = new(Config)
	config.tree, err = toml.LoadFile(*c)
	if err != nil {
		log.Panic(err)
	}
}

func DefaultConfig() Config {
	if config == nil {
		loadConfig()
	}
	return *config
}

func (c Config) GetSub(name string) Config {
	var b bool
	c.tree, b = c.tree.Get(name).(*toml.Tree)
	if b {
		return c
	}
	return Config{}
}
func (c Config) GetString(name string) string {
	v, b := c.tree.Get(name).(string)
	if b {
		return v
	}
	return ""
}

func (c Config) GetStringWithDefault(name string, def string) string {
	v, b := c.tree.Get(name).(string)
	if b {
		return v
	}
	return def
}
