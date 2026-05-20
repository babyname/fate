package seeddb

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/babyname/fate/resources"
)

func LoadEmbeddedCharacters() ([]SeedCharacter, error) {
	var seeds []SeedCharacter
	if err := json.Unmarshal(resources.CharacterJSON, &seeds); err != nil {
		return nil, fmt.Errorf("parse character.json: %w", err)
	}
	log.Printf("Loaded %d characters from embedded resources", len(seeds))
	return seeds, nil
}
