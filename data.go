package fate

import (
	"embed"
	"encoding/csv"
	"io/fs"
)

//下面的注释不能删，是用于向二进制文件中嵌入静态文件的
//go:embed data
var DataFiles embed.FS

func readData(f fs.File) ([][]string, error) {
	r := csv.NewReader(f)
	r.Comma = ','
	r.Comment = '#'

	// skip first line
	if _, err := r.Read(); err != nil {
		return [][]string{}, err
	}

	records, err := r.ReadAll()

	if err != nil {
		return [][]string{}, err
	}

	return records, nil
}
