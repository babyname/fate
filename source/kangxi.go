package source

import (
	"bufio"
	"io"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

type KangXi struct {
	CodePoint int
	Value     int
	Character string
	Strokes   int
}

func LoadKangXiChar(path string, hook func(kx KangXi) bool) error {
	//load polyphone from json file
	of, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer of.Close()
	//readline from open file
	br := bufio.NewReader(of)

	for {
		line, _, err := br.ReadLine()
		if err != nil {
			if err != io.EOF {
				slog.Error("readline error:", err)
			}
			break
		}
		if len(line) == 0 {
			continue
		}
		kx := decodeKangXi(string(line))
		if kx.CodePoint != 0 {
			if !hook(kx) {
				return nil
			}
		}
	}
	return nil
}

func decodeKangXi(line string) KangXi {
	chs := strings.Split(line, ",")
	if len(chs) != 4 {
		return KangXi{}
	}
	//fmt.Println("kangxi", chs[0][2:], "1:", chs[1], "2:", chs[2])
	cp, _ := strconv.ParseInt(chs[0][2:], 10, 32)
	v, _ := strconv.ParseInt(chs[1], 10, 32)
	stk, _ := strconv.ParseInt(chs[3], 10, 32)
	return KangXi{
		CodePoint: int(cp),
		Value:     int(v),
		Character: chs[2],
		Strokes:   int(stk),
	}
}
