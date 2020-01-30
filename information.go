package fate

import (
	"bufio"
	"fmt"
	"github.com/godcong/fate/config"
	"os"
)

type Information interface {
	Write(name Name)
	Finish()
}

func NewWithConfig(config config.Config) Information {
	return nil
}

type jsonInformation struct {
	path   string
	writer *bufio.Writer
}

func (j jsonInformation) Finish() {
}

func (j jsonInformation) Write(name Name) {

}

func jsonOutput(path string) Information {
	file, e := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_SYNC|os.O_RDWR, 0755)
	if e != nil {
		panic(fmt.Errorf("json output failed:%w", e))
	}
	writer := bufio.NewWriter(file)
	return &jsonInformation{
		path:   path,
		writer: writer,
	}
}
