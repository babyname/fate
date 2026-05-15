package api

import (
	"embed"
	"io/fs"
)

//go:embed static/*
var webFS embed.FS

func init() {
	var err error
	webSub, err = fs.Sub(webFS, "static")
	if err != nil {
		panic(err)
	}
}

var webSub fs.FS
