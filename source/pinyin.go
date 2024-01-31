package source

import (
	"bufio"
	"io"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

type PinYin struct {
	ID     int32
	Char   string
	Pinyin []string
}

func LoadPinYin(path string, hook func(yin *PinYin) bool) {
	py, err := os.Open(path)
	if err != nil {
		slog.Error("open file error:", err)
		return
	}
	defer py.Close()
	//readline from open file
	br := bufio.NewReader(py)

	for {
		line, _, err := br.ReadLine()
		if err != nil {
			if err != io.EOF {
				slog.Error("readline error:", err)
			}
			break
		}
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		pinyin := decodePinYin(string(line))
		if pinyin.ID != 0 {
			if !hook(pinyin) {
				return
			}
		}
	}
	return
}

func decodePinYin(s string) *PinYin {
	//decode from line U+idxxx: pinyin1,pinyin2  # char
	//get id
	idx := strings.Index(s, ":")
	id := s[2:idx]
	//slog.Debug("decode pinyin", "id", id)
	idInt, _ := strconv.ParseUint(id, 16, 32)
	//get pinyin
	if len(s) < idx+1 {
		return &PinYin{}
	}
	//slog.Debug("decode pinyin", "id", idInt)
	s = s[idx+1:]
	//slog.Debug("decode pinyin", "pinyin", s)
	idx = strings.Index(s, "#")
	py := s[:idx]
	//slog.Debug("decode pinyin", "pinyin", py)
	py = strings.TrimSpace(py)
	pinyin := strings.Split(py, ",")

	//get char
	if len(s) < idx+1 {
		return &PinYin{}
	}
	s = s[idx+1:]
	//slog.Debug("decode pinyin", "char", s)
	c := strings.TrimSpace(s)
	return &PinYin{
		ID:     int32(idInt),
		Char:   c,
		Pinyin: pinyin,
	}
}
