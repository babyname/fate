// Package analysis provides name analysis and report generation.
package analysis

import (
	"github.com/babyname/fate/ent"
	"github.com/babyname/fate/internal/rating"
	v2 "github.com/godcong/chronos/v2"
)

// calcScoreDetail calculates a detailed score breakdown for a name,
// combining culture impression, wu-xing/ba-zi, zodiac, and wu-ge scores.
func calcScoreDetail(nr *rating.NameRating, _ *ZhouYiResult, c1, c2 *ent.Character, fateData *v2.FateData) ScoreDetail {
	return ScoreDetail{
		WenHuaYinXiang: calcWenHuaScore(c1, c2),
		WuXingBaZi:     nr.WuXingScore,
		ShengXiao:      calcShengXiaoScore(c1, c2, fateData),
		WuGeShuLi:      nr.BiHuaScore,
	}
}

// calcWenHuaScore evaluates the cultural impression score of the given
// characters based on regularity, common usage level, and meaning.
func calcWenHuaScore(c1, c2 *ent.Character) float64 {
	score := 70.0
	if c1.Regular {
		score += 5
	}
	if c2.Regular {
		score += 5
	}
	if c1.CommonLevel > 0 && c1.CommonLevel <= 3 {
		score += 5
	}
	if c2.CommonLevel > 0 && c2.CommonLevel <= 3 {
		score += 5
	}
	if c1.Meaning != "" {
		score += 5
	}
	if c2.Meaning != "" {
		score += 5
	}
	if score > 100 {
		score = 100
	}
	return score
}

// calcShengXiaoScore evaluates the zodiac compatibility score for the
// given characters against the birth date's zodiac sign, considering
// wu-xing sheng (generating) and ke (overcoming) relationships.
func calcShengXiaoScore(c1, c2 *ent.Character, fateData *v2.FateData) float64 {
	score := 60.0
	if fateData != nil && fateData.Bazi != nil {
		zodiac := fateData.Bazi.Zodiac
		wx1 := c1.WuXing
		wx2 := c2.WuXing
		zodiacWuXing := getZodiacWuXing(zodiac)
		if zodiacWuXing != "" {
			if isSheng(zodiacWuXing, wx1) || isSheng(wx1, zodiacWuXing) {
				score += 10
			}
			if isSheng(zodiacWuXing, wx2) || isSheng(wx2, zodiacWuXing) {
				score += 10
			}
			if isKe(zodiacWuXing, wx1) {
				score -= 5
			}
			if isKe(zodiacWuXing, wx2) {
				score -= 5
			}
		}
	}
	if score > 100 {
		score = 100
	}
	if score < 0 {
		score = 0
	}
	return score
}
