package naming

import (
	"errors"

	"github.com/babyname/fate/v4/ent"
	"github.com/babyname/fate/v4/internal/analysis"
)

// ScoreName scores a known name using the same rater and result builder as generation.
func ScoreName(req GenerateRequest, first, second *ent.Character) (analysis.NameResult, error) {
	strokes := surnameStrokes(req.Filter, req.LastName)
	if strokes[0] <= 0 || (req.LastName[1] != nil && strokes[1] <= 0) {
		return analysis.NameResult{}, errors.New("surname is missing stroke counts")
	}
	if first.ScienceStroke <= 0 || second.ScienceStroke <= 0 {
		return analysis.NameResult{}, errors.New("given name is missing name-science stroke counts")
	}
	fateData, err := req.calculateFateData()
	if err != nil {
		return analysis.NameResult{}, err
	}
	rater := req.newRater(fateData, strokes)
	return analysis.BuildNameResultWithRater(0, surnameStr(req.LastName), first, second, strokes[0], strokes[1], fateData, rater), nil
}
