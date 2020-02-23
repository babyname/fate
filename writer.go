package fate

import (
	"bufio"
	"os"
)

// NewOutput ...
func NewOutput(path string) (*bufio.Writer, error) {
	file, e := os.Create(path)
	if e != nil {
		return nil, e
	}

	return bufio.NewWriter(file), nil

}
