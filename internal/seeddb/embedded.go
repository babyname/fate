package seeddb

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
)

//go:embed data/character.json
var embeddedSeedFS embed.FS

const embeddedSeedPath = "data/character.json"

// LoadEmbeddedCharacters loads the embedded character seed data.
// This data is compiled into the binary and is always available
// regardless of the working directory, ensuring the character
// database can always be populated.
func LoadEmbeddedCharacters() ([]SeedCharacter, error) {
	data, err := fs.ReadFile(embeddedSeedFS, embeddedSeedPath)
	if err != nil {
		return nil, fmt.Errorf("read embedded seed data: %w", err)
	}

	var seeds []SeedCharacter
	if err := json.Unmarshal(data, &seeds); err != nil {
		return nil, fmt.Errorf("parse embedded seed data: %w", err)
	}

	log.Printf("Loaded %d characters from embedded seed data", len(seeds))
	return seeds, nil
}
