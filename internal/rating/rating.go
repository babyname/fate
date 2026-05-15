package rating

import (
	"fmt"
	"math"
	"strings"

	"github.com/babyname/fate/ent"
	"github.com/babyname/fate/internal/wuge"
	"github.com/babyname/fate/internal/wuxing"
	v2 "github.com/godcong/chronos/v2"
)

type NameRating struct {
	WenHuaScore    float64
	WuXingScore    float64
	ShengXiaoScore float64
	WuGeScore      float64
	YinYunScore    float64
	TotalScore     float64
	Grade          string
	Interpret      string
	WenHuaDetail   string
	WuXingDetail   string
	ShengXiaoDetail string
	WuGeDetail     string
	YinYunDetail   string
}

type Rater struct {
	fateData *v2.FateData
	l1       int
	l2       int
}

func NewRater(fateData *v2.FateData) *Rater {
	return &Rater{fateData: fateData}
}

func NewRaterWithStrokes(fateData *v2.FateData, l1, l2 int) *Rater {
	return &Rater{fateData: fateData, l1: l1, l2: l2}
}

func (r *Rater) RateName(surname string, c1, c2 *ent.Character) *NameRating {
	rating := &NameRating{}

	r.rateWenHua(rating, c1, c2)
	r.rateWuXing(rating, c1, c2)
	r.rateShengXiao(rating, c1, c2)
	r.rateWuGe(rating, c1, c2)
	r.rateYinYun(rating, c1, c2)

	wenHuaWeight := 0.20
	wuXingWeight := 0.25
	shengXiaoWeight := 0.10
	wuGeWeight := 0.25
	yinYunWeight := 0.20

	rating.TotalScore = math.Round(
		(rating.WenHuaScore*wenHuaWeight+
			rating.WuXingScore*wuXingWeight+
			rating.ShengXiaoScore*shengXiaoWeight+
			rating.WuGeScore*wuGeWeight+
			rating.YinYunScore*yinYunWeight)*10,
	) / 10

	rating.Grade = scoreToGrade(rating.TotalScore)
	rating.Interpret = r.generateInterpret(surname, c1, c2, rating)

	return rating
}

func (r *Rater) rateWenHua(rating *NameRating, c1, c2 *ent.Character) {
	score := 60.0
	var details []string

	if c1.Regular {
		score += 5
		details = append(details, fmt.Sprintf("「%s」为常用字", c1.Char))
	}
	if c2.Regular {
		score += 5
		details = append(details, fmt.Sprintf("「%s」为常用字", c2.Char))
	}
	if c1.CommonLevel > 0 && c1.CommonLevel <= 3 {
		score += 5
		details = append(details, fmt.Sprintf("「%s」常用度高", c1.Char))
	}
	if c2.CommonLevel > 0 && c2.CommonLevel <= 3 {
		score += 5
		details = append(details, fmt.Sprintf("「%s」常用度高", c2.Char))
	}
	if c1.Meaning != "" {
		score += 5
		details = append(details, fmt.Sprintf("「%s」有明确释义", c1.Char))
	}
	if c2.Meaning != "" {
		score += 5
		details = append(details, fmt.Sprintf("「%s」有明确释义", c2.Char))
	}
	if c1.Nameable {
		score += 5
	}
	if c2.Nameable {
		score += 5
	}

	if score > 100 {
		score = 100
	}
	if score < 0 {
		score = 0
	}

	rating.WenHuaScore = score
	if len(details) > 0 {
		rating.WenHuaDetail = strings.Join(details, "; ")
	} else {
		rating.WenHuaDetail = "文化印象信息不足"
	}
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
			score -= 5
			details = append(details, fmt.Sprintf("「%s」五行属%s，为忌神，减分", c.Char, wx))
		} else {
			details = append(details, fmt.Sprintf("「%s」五行属%s，中性", c.Char, wx))
		}
	}

	wx1 := c1.WuXing
	wx2 := c2.WuXing
	if wx1 != "" && wx2 != "" {
		if isSheng(wx1, wx2) {
			score += 5
			details = append(details, fmt.Sprintf("五行相生：%s生%s", wx1, wx2))
		} else if isSheng(wx2, wx1) {
			score += 5
			details = append(details, fmt.Sprintf("五行相生：%s生%s", wx2, wx1))
		}
		if isKe(wx1, wx2) {
			score -= 5
			details = append(details, fmt.Sprintf("五行相克：%s克%s", wx1, wx2))
		} else if isKe(wx2, wx1) {
			score -= 5
			details = append(details, fmt.Sprintf("五行相克：%s克%s", wx2, wx1))
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

func (r *Rater) rateShengXiao(rating *NameRating, c1, c2 *ent.Character) {
	score := 60.0
	var details []string

	if r.fateData == nil || r.fateData.Bazi == nil {
		rating.ShengXiaoScore = score
		rating.ShengXiaoDetail = "无法获取生肖信息"
		return
	}

	zodiac := r.fateData.Bazi.Zodiac
	wx1 := c1.WuXing
	wx2 := c2.WuXing
	zodiacWuXing := getZodiacWuXing(zodiac)

	if zodiacWuXing != "" {
		details = append(details, fmt.Sprintf("生肖%s，五行属%s", zodiac, zodiacWuXing))

		if isSheng(zodiacWuXing, wx1) {
			score += 10
			details = append(details, fmt.Sprintf("「%s」(%s)为生肖所生，吉利", c1.Char, wx1))
		} else if isSheng(wx1, zodiacWuXing) {
			score += 8
			details = append(details, fmt.Sprintf("生肖为「%s」(%s)所生，吉利", c1.Char, wx1))
		}
		if isSheng(zodiacWuXing, wx2) {
			score += 10
			details = append(details, fmt.Sprintf("「%s」(%s)为生肖所生，吉利", c2.Char, wx2))
		} else if isSheng(wx2, zodiacWuXing) {
			score += 8
			details = append(details, fmt.Sprintf("生肖为「%s」(%s)所生，吉利", c2.Char, wx2))
		}
		if isKe(zodiacWuXing, wx1) {
			score -= 8
			details = append(details, fmt.Sprintf("生肖克「%s」(%s)，不利", c1.Char, wx1))
		}
		if isKe(zodiacWuXing, wx2) {
			score -= 8
			details = append(details, fmt.Sprintf("生肖克「%s」(%s)，不利", c2.Char, wx2))
		}
	}

	if score > 100 {
		score = 100
	}
	if score < 0 {
		score = 0
	}

	rating.ShengXiaoScore = score
	if len(details) > 0 {
		rating.ShengXiaoDetail = strings.Join(details, "; ")
	} else {
		rating.ShengXiaoDetail = "生肖信息不足"
	}
}

func (r *Rater) rateWuGe(rating *NameRating, c1, c2 *ent.Character) {
	if r.l1 == 0 {
		rating.WuGeScore = 60
		rating.WuGeDetail = "缺少姓氏笔画信息"
		return
	}

	f1 := c1.ScienceStroke
	f2 := c2.ScienceStroke
	ge := wuge.CalcWuGe(r.l1, r.l2, f1, f2)

	tianGe := ge.TianGe()
	renGe := ge.RenGe()
	diGe := ge.DiGe()
	waiGe := ge.WaiGe()
	zongGe := ge.ZongGe()

	tianDaYan := wuge.Find(tianGe)
	renDaYan := wuge.Find(renGe)
	diDaYan := wuge.Find(diGe)
	waiDaYan := wuge.Find(waiGe)
	zongDaYan := wuge.Find(zongGe)

	score := 40.0
	var details []string

	luckyCount := 0
	geItems := []struct {
		name  string
		daYan wuge.DaYan
	}{
		{"天格", tianDaYan},
		{"人格", renDaYan},
		{"地格", diDaYan},
		{"外格", waiDaYan},
		{"总格", zongDaYan},
	}

	for _, item := range geItems {
		switch {
		case strings.Contains(item.daYan.Lucky, "吉") && !strings.Contains(item.daYan.Lucky, "半"):
			luckyCount++
			details = append(details, fmt.Sprintf("%s%d(%s)吉", item.name, item.daYan.Number, item.daYan.Lucky))
		case strings.Contains(item.daYan.Lucky, "半吉"):
			details = append(details, fmt.Sprintf("%s%d(%s)", item.name, item.daYan.Number, item.daYan.Lucky))
		default:
			details = append(details, fmt.Sprintf("%s%d(%s)凶", item.name, item.daYan.Number, item.daYan.Lucky))
		}
	}

	score += float64(luckyCount) * 12.0

	if zongDaYan.Max {
		score += 5
		details = append(details, "总格为最大好运数")
	}

	sanCaiObj := wuxing.NewSanCai(tianGe, renGe, diGe)
	sanCaiStr := sanCaiObj.String()
	sanCaiLuck, _ := wuxing.GetWuXing(sanCaiStr)
	switch sanCaiLuck {
	case "大吉":
		score += 10
		details = append(details, fmt.Sprintf("三才%s大吉", sanCaiStr))
	case "吉", "吉多于凶":
		score += 7
		details = append(details, fmt.Sprintf("三才%s吉", sanCaiStr))
	case "中吉":
		score += 5
		details = append(details, fmt.Sprintf("三才%s中吉", sanCaiStr))
	case "凶多于吉", "吉凶参半":
		score -= 3
	case "凶", "大凶":
		score -= 8
		details = append(details, fmt.Sprintf("三才%s凶", sanCaiStr))
	}

	if score > 100 {
		score = 100
	}
	if score < 0 {
		score = 0
	}

	rating.WuGeScore = score
	rating.WuGeDetail = strings.Join(details, "; ")
}

func (r *Rater) rateYinYun(rating *NameRating, c1, c2 *ent.Character) {
	score := 60.0
	var details []string

	pinyin1 := c1.Pinyin
	pinyin2 := c2.Pinyin
	if len(pinyin1) == 0 || len(pinyin2) == 0 {
		rating.YinYunScore = score
		rating.YinYunDetail = "拼音信息不足"
		return
	}

	p1 := pinyin1[0]
	p2 := pinyin2[0]

	tone1 := getTone(p1)
	tone2 := getTone(p2)

	if tone1 != tone2 && tone1 != 0 && tone2 != 0 {
		score += 15
		details = append(details, "两字声调不同，抑扬顿挫")
	} else if tone1 == tone2 && tone1 != 0 {
		score -= 10
		details = append(details, "两字声调相同，略显平淡")
	}

	initial1 := getInitial(p1)
	initial2 := getInitial(p2)
	if initial1 == initial2 && initial1 != "" {
		score -= 8
		details = append(details, "声母相同，略有拗口")
	} else if initial1 != "" && initial2 != "" {
		score += 5
		details = append(details, "声母不同，发音清晰")
	}

	final1 := getFinal(p1)
	final2 := getFinal(p2)
	if final1 == final2 && final1 != "" {
		score -= 8
		details = append(details, "韵母相同，容易混淆")
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
		switch r.fateData.WuxingXiji.QiangRuo {
		case "强":
			parts = append(parts, "八字偏强，")
		case "弱":
			parts = append(parts, "八字偏弱，")
		default:
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

var shengMap = map[string]string{
	"木": "火", "火": "土", "土": "金", "金": "水", "水": "木",
}

var keMap = map[string]string{
	"木": "土", "土": "水", "水": "火", "火": "金", "金": "木",
}

func isSheng(a, b string) bool {
	if v, ok := shengMap[a]; ok {
		return v == b
	}
	return false
}

func isKe(a, b string) bool {
	if v, ok := keMap[a]; ok {
		return v == b
	}
	return false
}

func getZodiacWuXing(zodiac string) string {
	m := map[string]string{
		"鼠": "水", "牛": "土", "虎": "木", "兔": "木",
		"龙": "土", "蛇": "火", "马": "火", "羊": "土",
		"猴": "金", "鸡": "金", "狗": "土", "猪": "水",
	}
	if v, ok := m[zodiac]; ok {
		return v
	}
	return ""
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

func getFinal(pinyin string) string {
	runes := []rune(pinyin)
	if len(runes) <= 1 {
		return ""
	}
	return string(runes[1:])
}
