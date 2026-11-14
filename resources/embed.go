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

// StaticSub is the only embedded filesystem. Embedding the whole resources
// directory would pull temporary backups into every compiled binary.
