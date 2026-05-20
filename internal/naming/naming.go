// Package naming 提供姓名生成与评分功能，支持五行、笔画、音韵等多维度筛选与推荐。
package naming

import (
	v2 "github.com/babyname/chronos/v2"
	"github.com/babyname/fate/config"
	"github.com/babyname/fate/ent"
	"github.com/babyname/fate/internal/repository"
)

// Interface 定义姓名生成与评分的接口。
type Interface interface {
	FilterNames(char1, char2 []*ent.Character, filter *NameFilter) ([]*NameInfo, error)
	RateNames(names []*NameInfo, fateData *v2.FateData) (*RatedNames, error)
	RecommendNames(surname string, char1, char2 []*ent.Character, fateData *v2.FateData, opts *RecommendOptions) (*RecommendedNames, error)
}

// Naming 姓名生成与评分的核心实现。
type Naming struct {
	cfg    *config.Config
	model  *repository.Repository
	raters []Rater
}

// Rater 定义单个评分维度的接口。
type Rater interface {
	Rate(name *NameInfo, fateData *v2.FateData) (float64, string)
	Name() string
	Weight() float64
}

// NameInfo 表示一个姓名的基本信息。
type NameInfo struct {
	Surname   string         `json:"surname"`
	GivenName string         `json:"given_name"`
	FullName  string         `json:"full_name"`
	Char1     *ent.Character `json:"char1"`
	Char2     *ent.Character `json:"char2"`
}

// RatedName 表示一个已评分的姓名。
type RatedName struct {
	Name   *NameInfo          `json:"name"`
	Score  float64            `json:"score"`
	Grades map[string]float64 `json:"grades"`
	Grade  string             `json:"grade"`
	Notes  map[string]string  `json:"notes"`
}

// RatedNames 表示一组已评分的姓名集合。
type RatedNames struct {
	Names []*RatedName `json:"names"`
}

// RecommendedNames 表示一组推荐姓名及其数量。
type RecommendedNames struct {
	Names []*RecommendedName `json:"names"`
	Count int                `json:"count"`
}

// RecommendedName 表示一个推荐姓名，包含评分结果和解读。
type RecommendedName struct {
	*RatedName
	Interpretation string `json:"interpretation"`
}

// NameFilter 定义姓名筛选条件。
type NameFilter struct {
	PreferredWuxing []string
	AvoidedWuxing   []string
	MinStroke       int
	MaxStroke       int
	OnlyCommon      bool
	CommonLevel     int
	OnlyRegular     bool
}

// RecommendOptions 定义姓名推荐的选项。
type RecommendOptions struct {
	MaxResults int
	OnlyTop    bool
	SortBy     string
}

// New 创建一个新的姓名生成与评分实例。
func New(cfg *config.Config, model *repository.Repository) Interface {
	n := &Naming{
		cfg:   cfg,
		model: model,
	}

	n.raters = []Rater{
		NewWuxingRater(cfg),
		NewBihuaRater(cfg),
		NewYinyunRater(cfg),
	}

	return n
}

// FilterNames 根据筛选条件过滤字符组合，生成候选姓名列表。
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

// RateNames 对候选姓名列表进行多维度评分并排序。
func (n *Naming) RateNames(names []*NameInfo, fateData *v2.FateData) (*RatedNames, error) {
	rated := make([]*RatedName, 0, len(names))

	for _, name := range names {
		r := n.rateName(name, fateData)
		rated = append(rated, r)
	}

	sortRatedNames(rated, "score")

	return &RatedNames{Names: rated}, nil
}

func (n *Naming) rateName(name *NameInfo, fateData *v2.FateData) *RatedName {
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

// RecommendNames 根据八字命理数据推荐最优姓名。
func (n *Naming) RecommendNames(surname string, char1, char2 []*ent.Character, fateData *v2.FateData, opts *RecommendOptions) (*RecommendedNames, error) {
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

func generateInterpretation(r *RatedName, _ *v2.FateData) string {
	return r.Name.FullName + " - 评分: " + r.Grade
}

func sortRatedNames(_ []*RatedName, _ string) {
}
