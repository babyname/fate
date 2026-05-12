package rating

import (
	"fmt"
	"math"
	"strings"

	v2 "github.com/godcong/chronos/v2"
	"github.com/babyname/fate/ent"
)

type NameRating struct {
	WuXingScore  float64
	BiHuaScore   float64
	YinYunScore  float64
	TotalScore   float64
	Grade        string
	Interpret    string
	WuXingDetail string
	BiHuaDetail  string
	YinYunDetail string
}

type Rater struct {
	fateData *v2.FateData
}

func NewRater(fateData *v2.FateData) *Rater {
	return &Rater{fateData: fateData}
}

func (r *Rater) RateName(surname string, c1, c2 *ent.Character) *NameRating {
	rating := &NameRating{}

	r.rateWuXing(rating, c1, c2)
	r.rateBiHua(rating, c1, c2)
	r.rateYinYun(rating, c1, c2)

	wuxingWeight := 0.4
	bihuaWeight := 0.4
	yinyunWeight := 0.2

	rating.TotalScore = math.Round(
		(rating.WuXingScore*wuxingWeight+
			rating.BiHuaScore*bihuaWeight+
			rating.YinYunScore*yinyunWeight)*10,
	) / 10

	rating.Grade = scoreToGrade(rating.TotalScore)
	rating.Interpret = r.generateInterpret(surname, c1, c2, rating)

	return rating
}

func (r *Rater) rateWuXing(rating *NameRating, c1, c2 *ent.Character) {
	if r.fateData == nil || r.fateData.WuxingXiji == nil {
		rating.WuXingScore = 60
		rating.WuXingDetail = "无法获取五行信息"
		return
	}

	xiWuxing := r.fateData.WuxingXiji.XiWuxing
	jiWuxing := r.fateData.WuxingXiji.JiWuxing

	score := 50.0
	var details []string

	chars := []*ent.Character{c1, c2}
	for _, c := range chars {
		wx := c.WuXing
		if wx == "" {
			continue
		}

		isXi := contains(xiWuxing, wx)
		isJi := contains(jiWuxing, wx)

		if isXi {
			score += 15
			details = append(details, fmt.Sprintf("「%s」五行属%s，为喜用神，加分", c.Char, wx))
		} else if isJi {
			details = append(details, fmt.Sprintf("「%s」五行属%s，为忌神，减分", c.Char, wx))
		} else {
			details = append(details, fmt.Sprintf("「%s」五行属%s，中性", c.Char, wx))
		}
	}

	if score > 100 {
		score = 100
	}
	if score < 0 {
		score = 0
	}

	rating.WuXingScore = score
	if len(details) > 0 {
		rating.WuXingDetail = strings.Join(details, "; ")
	} else {
		rating.WuXingDetail = "五行信息不足"
	}
}

func (r *Rater) rateBiHua(rating *NameRating, c1, c2 *ent.Character) {
	score := 70.0
	var details []string

	chars := []*ent.Character{c1, c2}
	for _, c := range chars {
		stroke := c.ScienceStroke
		if stroke >= 8 && stroke <= 15 {
			score += 5
			details = append(details, fmt.Sprintf("「%s」笔画%d，在吉利范围内", c.Char, stroke))
		} else if stroke < 5 || stroke > 20 {
			details = append(details, fmt.Sprintf("「%s」笔画%d，偏大或偏小", c.Char, stroke))
		}
	}

	if c1.ScienceStroke == c2.ScienceStroke {
		score += 5
		details = append(details, "两字笔画相同，书写匀称")
	}

	if score > 100 {
		score = 100
	}
	if score < 0 {
		score = 0
	}

	rating.BiHuaScore = score
	if len(details) > 0 {
		rating.BiHuaDetail = strings.Join(details, "; ")
	} else {
		rating.BiHuaDetail = "笔画信息不足"
	}
}

func (r *Rater) rateYinYun(rating *NameRating, c1, c2 *ent.Character) {
	score := 70.0
	var details []string

	pinyin1 := c1.Pinyin
	pinyin2 := c2.Pinyin
	if len(pinyin1) == 0 || len(pinyin2) == 0 {
		rating.YinYunScore = 60
		rating.YinYunDetail = "拼音信息不足"
		return
	}

	p1 := pinyin1[0]
	p2 := pinyin2[0]

	tone1 := getTone(p1)
	tone2 := getTone(p2)

	if tone1 != tone2 {
		score += 15
		details = append(details, "两字声调不同，抑扬顿挫")
	} else {
		score -= 10
		details = append(details, "两字声调相同，略显平淡")
	}

	initial1 := getInitial(p1)
	initial2 := getInitial(p2)
	if initial1 == initial2 {
		score -= 5
		details = append(details, "声母相同，略有拗口")
	}

	if score > 100 {
		score = 100
	}
	if score < 0 {
		score = 0
	}

	rating.YinYunScore = score
	if len(details) > 0 {
		rating.YinYunDetail = strings.Join(details, "; ")
	} else {
		rating.YinYunDetail = "音韵信息不足"
	}
}

func (r *Rater) generateInterpret(surname string, c1, c2 *ent.Character, rating *NameRating) string {
	var parts []string

	parts = append(parts, fmt.Sprintf("「%s%s」", surname, c1.Char+c2.Char))

	if r.fateData != nil && r.fateData.WuxingXiji != nil {
		parts = append(parts, fmt.Sprintf("日主%s五行属%s，", r.fateData.WuxingXiji.DayGan, r.fateData.WuxingXiji.DayWuxing))
		if r.fateData.WuxingXiji.QiangRuo == "强" {
			parts = append(parts, "八字偏强，")
		} else if r.fateData.WuxingXiji.QiangRuo == "弱" {
			parts = append(parts, "八字偏弱，")
		} else {
			parts = append(parts, "八字中和，")
		}
		parts = append(parts, fmt.Sprintf("喜用%s。", strings.Join(r.fateData.WuxingXiji.XiWuxing, "、")))
	}

	wx1 := c1.WuXing
	wx2 := c2.WuXing
	if wx1 != "" && wx2 != "" {
		parts = append(parts, fmt.Sprintf("名字五行组合：%s%s", wx1, wx2))
		if r.fateData != nil && r.fateData.WuxingXiji != nil {
			if contains(r.fateData.WuxingXiji.XiWuxing, wx1) || contains(r.fateData.WuxingXiji.XiWuxing, wx2) {
				parts = append(parts, "与八字喜用神相合。")
			}
		}
	}

	parts = append(parts, fmt.Sprintf("综合评分%.1f分（%s）", rating.TotalScore, rating.Grade))

	return strings.Join(parts, "")
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

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func getTone(pinyin string) int {
	runes := []rune(pinyin)
	if len(runes) == 0 {
		return 0
	}
	last := runes[len(runes)-1]
	if last >= '1' && last <= '4' {
		return int(last - '0')
	}
	return 0
}

func getInitial(pinyin string) string {
	runes := []rune(pinyin)
	if len(runes) == 0 {
		return ""
	}
	return string(runes[0])
}
