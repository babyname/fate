package resources

import (
	"embed"
	"io/fs"
)

//go:embed character.json
var CharacterJSON []byte

//go:embed all:static
var staticFS embed.FS

var StaticSub fs.FS

func init() {
	var err error
	StaticSub, err = fs.Sub(staticFS, "static")
	if err != nil {
		panic(err)
	}
}

func OpenFS() fs.FS {
	return EmbeddedFS
}

//go:embed *
var EmbeddedFS embed.FS
