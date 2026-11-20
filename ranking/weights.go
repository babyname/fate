// Package ranking combines traditional and sound-and-meaning scores.
package ranking

import (
	"errors"
	"math"
)

// Weight holds the traditional percentage. Sound and meaning receive the rest.
// Its zero value assigns the whole score to sound and meaning.
type Weight struct {
	traditional int
}

// NewWeight accepts a whole-number traditional percentage from 0 to 100.
func NewWeight(traditional int) (Weight, error) {
	if traditional < 0 || traditional > 100 {
		return Weight{}, errors.New("traditional_weight must be between 0 and 100")
	}
	return Weight{traditional: traditional}, nil
}

// Result holds the two component scores and their weighted total.
type Result struct {
	Total             float64 `json:"total"`
	Traditional       float64 `json:"traditional"`
	SoundAndMeaning   float64 `json:"sound_and_meaning"`
	TraditionalWeight int     `json:"traditional_weight"`
}

// Field identifies a score column used to order a candidate catalogue.
type Field string

const (
	TotalScore        Field = "total"
	TraditionalScore  Field = "traditional"
	SoundMeaningScore Field = "sound_and_meaning"
)

// Combine rounds each component to tenths, then rounds the weighted total.
// Integer tenths keep the total consistent with the component scores shown to users.
func (w Weight) Combine(traditional, soundAndMeaning float64) Result {
	traditionalTenths := int(math.Round(traditional * 10))
	soundTenths := int(math.Round(soundAndMeaning * 10))
	weightedTenths := traditionalTenths*w.traditional + soundTenths*(100-w.traditional)
	return Result{
		Total:             math.Round(float64(weightedTenths)/100) / 10,
		Traditional:       float64(traditionalTenths) / 10,
		SoundAndMeaning:   float64(soundTenths) / 10,
		TraditionalWeight: w.traditional,
	}
}
