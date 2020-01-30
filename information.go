package fate

import (
	"fmt"
	"github.com/godcong/fate/config"
	"github.com/goextension/log"
	"os"
)

type Information interface {
	Write(name Name)
}

func NewWithConfig(config config.Config) Information {
	return nil
}

type jsonInformation struct {
	path string
	file *os.File
}

func (j jsonInformation) Write(name Name) {
	log.Info(name)
}

func jsonOutput(path string) Information {
	file, e := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_SYNC|os.O_RDWR, 0755)
	if e != nil {
		panic(fmt.Errorf("json output failed:%w", e))
	}
	return &jsonInformation{
		path: path,
		file: file,
	}
}
