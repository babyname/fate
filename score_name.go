package fate

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/babyname/fate/v4/internal/naming"
	"github.com/babyname/fate/v4/internal/repository"
	namesvc "github.com/babyname/fate/v4/internal/service/naming"
	"github.com/babyname/fate/v4/ranking"
)

// ScoreNameRequest contains the details for evaluating an existing name.
type ScoreNameRequest struct {
	LastName  string
	FirstName string
	Born      time.Time
	Sex       Sex
	// FilterType selects the surname stroke convention used for generation.
	FilterType CharacterFilterType
	// XiYongMethod accepts the same values as FilterOption.XiYongMethod.
	XiYongMethod string
	// RecommendationWeight is optional. Nil uses the original five scores.
	RecommendationWeight *ranking.Weight
}

// ScoreName evaluates a one- or two-character surname and a two-character given
// name using the engine's character database. The result's Rank is zero because
// the name is evaluated on its own.
func ScoreName(ctx context.Context, engine Fate, req ScoreNameRequest) (NameResult, error) {
	lastName := strings.TrimSpace(req.LastName)
	givenName := []rune(strings.TrimSpace(req.FirstName))
	if !isHanName([]rune(lastName), 1, 2) {
		return NameResult{}, errors.New("surname must contain one or two Chinese characters")
	}
	if !isHanName(givenName, 2, 2) {
		return NameResult{}, errors.New("given name must contain two Chinese characters")
	}
	if req.Born.IsZero() {
		return NameResult{}, errors.New("birth time is required")
	}
	if req.Sex != SexBoy && req.Sex != SexGirl {
		return NameResult{}, errors.New("invalid sex")
	}

	repo := engine.Repo()
	surname, err := repo.QueryLastName(ctx, splitSurname(lastName))
	if err != nil {
		return NameResult{}, err
	}
	first, err := repo.GetCharacter(ctx, repository.Char(string(givenName[0])))
	if err != nil {
		return NameResult{}, fmt.Errorf("look up given name character %c: %w", givenName[0], err)
	}
	second := first
	if givenName[0] != givenName[1] {
		second, err = repo.GetCharacter(ctx, repository.Char(string(givenName[1])))
		if err != nil {
			return NameResult{}, fmt.Errorf("look up given name character %c: %w", givenName[1], err)
		}
	}

	flt := DefaultFilterOption()
	flt.CharacterFilterType = req.FilterType
	flt.XiYongMethod = req.XiYongMethod
	return namesvc.ScoreName(namesvc.GenerateRequest{
		LastName:             surname,
		Born:                 req.Born,
		Sex:                  naming.Sex(req.Sex),
		Filter:               NewFilter(flt),
		RecommendationWeight: req.RecommendationWeight,
	}, first, second)
}

func isHanName(chars []rune, minLength, maxLength int) bool {
	if len(chars) < minLength || len(chars) > maxLength {
		return false
	}
	for _, char := range chars {
		if !unicode.Is(unicode.Han, char) {
			return false
		}
	}
	return true
}
