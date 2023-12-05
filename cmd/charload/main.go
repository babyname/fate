package main

import (
	"encoding/json"
	"log/slog"
	"os"

	"github.com/babyname/fate/ent"
)

var path = "./char.json"

func main() {

	file, err := os.ReadFile(path)
	if err != nil {
		slog.Error("read file error", "info", err)
		return
	}
	var char map[any]ent.NCharacter
	err = json.Unmarshal(file, &char)
	if err != nil {
		slog.Error("unmarshal error", "info", err)
		return
	}
	//todo
}
