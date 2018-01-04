package config

import (
	"flag"
	"io/ioutil"
	"log"

	"github.com/pelletier/go-toml"
)

type Config struct {
	Database []map[string]string `toml:"database"`
	//Override []map[string]string `toml:"override"`
}

func LoadConfig() {
	c := flag.String("c", "config.toml", "load an toml config file")
	bs, err := ioutil.ReadFile(*c)
	log.Println((string(bs)))
	if err != nil {
		log.Println(err)
	}

	cfg := Config{}
	e := toml.Unmarshal(bs, &cfg)
	log.Println(cfg, e)
}
