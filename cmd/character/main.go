package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

func loadCharFile(path string) ([]byte, error) {
	opened, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(opened)
}

func main() {
	path := ""
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	if path == "" {
		return
	}
	fmt.Println("Load", path)
	file, err := loadCharFile(path)
	if err != nil {
		return
	}

	codes := []rune(string(file))
	filter := map[rune]bool{}
	for _, code := range codes {
		if strings.TrimSpace(string(code)) == "" {
			continue
		}
		if filterRune(code) {
			continue
		}
		filter[code] = true
	}
	var data []rune
	for r, b := range filter {
		if b {
			data = append(data, r)
		}
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i] < data[j]
	})

	if err := writeDataToFile(path+".new", string(data)); err != nil {
		return
	}

}

var filterRunes = []rune{
	'、',
}

func filterRune(r rune) bool {
	for _, filterRune := range filterRunes {
		if filterRune == r {
			return true
		}
	}
	return false
}

func writeDataToFile(path string, data string) error {
	return ioutil.WriteFile(path, []byte(data), 0755)
}
