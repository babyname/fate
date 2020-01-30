package fate

import (
	"fmt"
	"github.com/godcong/fate/config"
	"os"
)

type Information interface {
	Write(src string)
}

func NewWithConfig(config config.Config) Information {
	return nil
}

type jsonInformation struct {
	path string
	file *os.File
}

func (j jsonInformation) Write(src string) {
	panic("implement me")
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
