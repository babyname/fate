package resources

import (
	"bytes"
	"compress/gzip"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
)

//go:embed character.json
var CharacterJSON []byte

//go:embed fate.db.gz
var FateDBGZ []byte

//go:embed *
var EmbeddedFS embed.FS

func HasDB() bool {
	return len(FateDBGZ) > 0
}

func FateDBGZSize() int {
	return len(FateDBGZ)
}

func ExtractDB(destPath string) error {
	gr, err := gzip.NewReader(bytes.NewReader(FateDBGZ))
	if err != nil {
		return fmt.Errorf("create gzip reader: %w", err)
	}
	defer gr.Close()

	out, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("create db file: %w", err)
	}
	defer out.Close()

	written, err := io.Copy(out, gr)
	if err != nil {
		os.Remove(destPath)
		return fmt.Errorf("decompress db: %w", err)
	}
	log.Printf("[RES] Extracted database to %s: %d bytes", destPath, written)
	return nil
}

func OpenFS() fs.FS {
	return EmbeddedFS
}
