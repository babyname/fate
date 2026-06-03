package naming

import (
	"context"

	"github.com/babyname/fate/v4/internal/chronosfate"
	"github.com/babyname/fate/v4/config"
	"github.com/babyname/fate/v4/ent"
	"github.com/babyname/fate/v4/internal/repository"
)

type Interface interface {
	FilterNames(char1, char2 []*ent.Character, filter *NameFilter) ([]*NameInfo, error)
	RateNames(names []*NameInfo, fateData *chronosfate.FateData) (*RatedNames, error)
	RecommendNames(surname string, char1, char2 []*ent.Character, fateData *chronosfate.FateData, opts *RecommendOptions) (*RecommendedNames, error)
}

type Naming struct {
	cfg    *config.Config
	model  *repository.Repository
	raters []Rater
}

type Rater interface {
	Rate(name *NameInfo, fateData *chronosfate.FateData) (float64, string)
	Name() string
	Weight() float64
}

type NameInfo struct {
	Surname   string         `json:"surname"`
	GivenName string         `json:"given_name"`
	FullName  string         `json:"full_name"`
	Char1     *ent.Character `json:"char1"`
	Char2     *ent.Character `json:"char2"`
}

type RatedName struct {
	Name   *NameInfo          `json:"name"`
	Score  float64            `json:"score"`
	Grades map[string]float64 `json:"grades"`
	Grade  string             `json:"grade"`
	Notes  map[string]string  `json:"notes"`
}

type RatedNames struct {
	Names []*RatedName `json:"names"`
}

type RecommendedNames struct {
	Names []*RecommendedName `json:"names"`
	Count int                `json:"count"`
}

type RecommendedName struct {
	*RatedName
	Interpretation string `json:"interpretation"`
}

type NameFilter struct {
	PreferredWuxing []string
	AvoidedWuxing   []string
	MinStroke       int
	MaxStroke       int
	OnlyCommon      bool
	CommonLevel     int
	OnlyRegular     bool
}

type RecommendOptions struct {
	MaxResults int
	OnlyTop    bool
	SortBy     string
}

func New(cfg *config.Config, model *repository.Repository) Interface {
	n := &Naming{
		cfg:   cfg,
		model: model,
	}

	poetryChars := make(map[string]bool)
	if model != nil {
		if chars, err := model.QueryPoetryChars(context.Background()); err == nil {
			for _, ch := range chars {
				poetryChars[ch] = true
			}
		}
	}

	n.raters = []Rater{
		NewWuxingRater(cfg),
		NewBihuaRater(cfg),
		NewYinyunRater(cfg),
		NewPoetryRater(cfg, poetryChars),
	}

	return n
}

func (n *Naming) FilterNames(char1, char2 []*ent.Character, filter *NameFilter) ([]*NameInfo, error) {
	if filter == nil {
		filter = &NameFilter{
			MinStroke:   n.cfg.Filter.MinStroke,
			MaxStroke:   n.cfg.Filter.MaxStroke,
			OnlyCommon:  n.cfg.Filter.EnableCommonFilter,
			CommonLevel: n.cfg.Filter.CommonLevel,
			OnlyRegular: n.cfg.Filter.EnableRegularFilter,
		}
	}

	var filtered []*NameInfo

	for _, c1 := range char1 {
		for _, c2 := range char2 {
			if !n.filterChar(c1, filter) || !n.filterChar(c2, filter) {
				continue
			}

			info := &NameInfo{
				GivenName: c1.Char + c2.Char,
				Char1:     c1,
				Char2:     c2,
			}
			filtered = append(filtered, info)
		}
	}

	return filtered, nil
}

func (n *Naming) filterChar(c *ent.Character, filter *NameFilter) bool {
	stroke := c.KangxiStroke
	if stroke < filter.MinStroke || stroke > filter.MaxStroke {
		return false
	}

	if filter.OnlyRegular && !c.Regular {
		return false
	}

	return true
}

func (n *Naming) RateNames(names []*NameInfo, fateData *chronosfate.FateData) (*RatedNames, error) {
	rated := make([]*RatedName, 0, len(names))

	for _, name := range names {
		r := n.rateName(name, fateData)
		rated = append(rated, r)
	}

	sortRatedNames(rated, "score")

	return &RatedNames{Names: rated}, nil
}

func (n *Naming) rateName(name *NameInfo, fateData *chronosfate.FateData) *RatedName {
	grades := make(map[string]float64)
	notes := make(map[string]string)
	var totalScore float64

	for _, rater := range n.raters {
		score, note := rater.Rate(name, fateData)
		grades[rater.Name()] = score
		notes[rater.Name()] = note
		totalScore += score * rater.Weight()
	}

	return &RatedName{
		Name:   name,
		Score:  totalScore,
		Grades: grades,
		Grade:  scoreToGrade(totalScore),
		Notes:  notes,
	}
}

func (n *Naming) RecommendNames(surname string, char1, char2 []*ent.Character, fateData *chronosfate.FateData, opts *RecommendOptions) (*RecommendedNames, error) {
	if opts == nil {
		opts = &RecommendOptions{
			MaxResults: n.cfg.Output.MaxResults,
			SortBy:     "score",
		}
	}

	filtered, err := n.FilterNames(char1, char2, nil)
	if err != nil {
		return nil, err
	}

	for _, name := range filtered {
		name.Surname = surname
		name.FullName = surname + name.GivenName
	}

	rated, err := n.RateNames(filtered, fateData)
	if err != nil {
		return nil, err
	}

	sortRatedNames(rated.Names, opts.SortBy)

	if opts.MaxResults > 0 && len(rated.Names) > opts.MaxResults {
		rated.Names = rated.Names[:opts.MaxResults]
	}

	rec := make([]*RecommendedName, 0, len(rated.Names))
	for _, r := range rated.Names {
		rec = append(rec, &RecommendedName{
			RatedName:      r,
			Interpretation: generateInterpretation(r, fateData),
		})
	}

	return &RecommendedNames{
		Names: rec,
		Count: len(rec),
	}, nil
}

func scoreToGrade(score float64) string {
	switch {
	case score >= 90:
		return "上上"
	case score >= 80:
		return "上吉"
	case score >= 70:
		return "中吉"
	case score >= 60:
		return "中平"
	case score >= 50:
		return "中下"
	default:
		return "下下"
	}
}

func generateInterpretation(r *RatedName, _ *chronosfate.FateData) string {
	return r.Name.FullName + " - 评分: " + r.Grade
}

func sortRatedNames(_ []*RatedName, _ string) {
}
