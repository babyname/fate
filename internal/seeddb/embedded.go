package seeddb

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"
)

//go:embed data/character.json
var embeddedSeedFS embed.FS

func LoadEmbeddedCharacters() ([]SeedCharacter, error) {
	data, err := embeddedSeedFS.ReadFile("data/character.json")
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
