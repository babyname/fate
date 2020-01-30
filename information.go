package fate

import (
	"bufio"
	"fmt"
	"github.com/godcong/fate/config"
	"os"
)

type Information interface {
	Write(name ...Name) error
	Finish() error
}

func NewWithConfig(config config.Config) Information {
	return nil
}

type jsonInformation struct {
	path string
	file *os.File
}

func (j *jsonInformation) Finish() error {
	return j.file.Close()
}

func (j *jsonInformation) Write(names ...Name) error {
	w := bufio.NewWriter(j.file)
	for _, n := range names {
		_, _ = w.WriteString(n.String())
		_, _ = w.WriteString("\r\n")
	}
	return w.Flush()

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
