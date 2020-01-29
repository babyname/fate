package information

import "os"

type jsonInformation struct {
	path string
	file *os.File
}

func jsonOutput(path string) Information {
	file, e := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_SYNC|os.O_RDWR, 0755)
	if e != nil {
		return nil
	}
	return &jsonInformation{
		path: path,
		file: file,
	}
}
