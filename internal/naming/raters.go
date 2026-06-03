package naming

import (
	"github.com/babyname/fate/v4/internal/chronosfate"
	"github.com/babyname/fate/v4/config"
)

type WuxingRater struct {
	cfg *config.Config
}

func NewWuxingRater(cfg *config.Config) *WuxingRater {
	return &WuxingRater{cfg: cfg}
}

func (r *WuxingRater) Name() string {
	return "wuxing"
}

func (r *WuxingRater) Weight() float64 {
	return r.cfg.Rate.WuxingWeight
}

func (r *WuxingRater) Rate(name *NameInfo, fateData *chronosfate.FateData) (float64, string) {
	if fateData == nil {
		return 60.0, "缺少八字信息"
	}

	score := 50.0
	var notes []string

	wuxing1 := name.Char1.WuXing
	wuxing2 := name.Char2.WuXing

	xi := fateData.WuxingXijiInfo.Xi
	ji := fateData.WuxingXijiInfo.Ji

	if wuxing1 == xi {
		score += 20.0
		notes = append(notes, name.Char1.Char+"("+wuxing1+")符合喜用五行")
	}
	if wuxing2 == xi {
		score += 20.0
		notes = append(notes, name.Char2.Char+"("+wuxing2+")符合喜用五行")
	}

	if wuxing1 == ji {
		score -= 15.0
		notes = append(notes, name.Char1.Char+"("+wuxing1+")为忌神五行")
	}
	if wuxing2 == ji {
		score -= 15.0
		notes = append(notes, name.Char2.Char+"("+wuxing2+")为忌神五行")
	}

	if isWuXingSheng(wuxing1, wuxing2) {
		score += 10.0
		notes = append(notes, wuxing1+"生"+wuxing2+"，五行相生")
	}

	if score > 100.0 {
		score = 100.0
	}
	if score < 0.0 {
		score = 0.0
	}

	var noteStr string
	if len(notes) > 0 {
		noteStr = joinNotes(notes)
	} else {
		noteStr = "五行中和"
	}

	return score, noteStr
}

type BihuaRater struct {
	cfg *config.Config
}

func NewBihuaRater(cfg *config.Config) *BihuaRater {
	return &BihuaRater{cfg: cfg}
}

func (r *BihuaRater) Name() string {
	return "bihua"
}

func (r *BihuaRater) Weight() float64 {
	return r.cfg.Rate.BihuaWeight
}

func (r *BihuaRater) Rate(name *NameInfo, _ *chronosfate.FateData) (float64, string) {
	score := 70.0
	var notes []string

	stroke1 := name.Char1.KangxiStroke
	stroke2 := name.Char2.KangxiStroke

	if stroke1 >= 5 && stroke1 <= 20 {
		score += 10.0
	} else {
		score -= 10.0
		notes = append(notes, name.Char1.Char+"笔画偏多/偏少")
	}
	if stroke2 >= 5 && stroke2 <= 20 {
		score += 10.0
	} else {
		score -= 10.0
		notes = append(notes, name.Char2.Char+"笔画偏多/偏少")
	}

	diff := abs(stroke1-stroke2)
	if diff <= 5 {
		score += 10.0
		notes = append(notes, "笔画搭配匀称")
	} else if diff > 10 {
		score -= 10.0
		notes = append(notes, "笔画搭配不均")
	}

	if isLuckyStroke(stroke1) {
		score += 5.0
	}
	if isLuckyStroke(stroke2) {
		score += 5.0
	}

	if score > 100.0 {
		score = 100.0
	}
	if score < 0.0 {
		score = 0.0
	}

	var noteStr string
	if len(notes) > 0 {
		noteStr = joinNotes(notes)
	} else {
		noteStr = "笔画适中"
	}

	return score, noteStr
}

type YinyunRater struct {
	cfg *config.Config
}

func NewYinyunRater(cfg *config.Config) *YinyunRater {
	return &YinyunRater{cfg: cfg}
}

func (r *YinyunRater) Name() string {
	return "yinyun"
}

func (r *YinyunRater) Weight() float64 {
	return r.cfg.Rate.YinyunWeight
}

func (r *YinyunRater) Rate(name *NameInfo, _ *chronosfate.FateData) (float64, string) {
	score := 70.0
	var notes []string

	pinyin1 := firstPinyin(name.Char1.Pinyin)
	pinyin2 := firstPinyin(name.Char2.Pinyin)

	if pinyin1 != "" && pinyin2 != "" {
		tone1 := getTone(pinyin1)
		tone2 := getTone(pinyin2)

		if tone1 != tone2 && tone1 != 0 && tone2 != 0 {
			score += 15.0
			notes = append(notes, "声调搭配抑扬顿挫")
		} else {
			score -= 10.0
			notes = append(notes, "声调相同，略显平淡")
		}

		sheng1 := getShengMu(pinyin1)
		sheng2 := getShengMu(pinyin2)

		if sheng1 != sheng2 && sheng1 != "" && sheng2 != "" {
			score += 10.0
			notes = append(notes, "声母不同，发音清晰")
		} else {
			score -= 5.0
			notes = append(notes, "声母相同，易绕口")
		}

		yun1 := getYunMu(pinyin1)
		yun2 := getYunMu(pinyin2)

		if yun1 != yun2 && yun1 != "" && yun2 != "" {
			score += 5.0
		}
	} else {
		score -= 20.0
		notes = append(notes, "拼音信息不足")
	}

	if score > 100.0 {
		score = 100.0
	}
	if score < 0.0 {
		score = 0.0
	}

	var noteStr string
	if len(notes) > 0 {
		noteStr = joinNotes(notes)
	} else {
		noteStr = "音律中和"
	}

	return score, noteStr
}

func isWuXingSheng(wx1, wx2 string) bool {
	sheng := map[string]string{
		"木": "火", "火": "土", "土": "金", "金": "水", "水": "木",
	}
	if next, ok := sheng[wx1]; ok {
		return next == wx2
	}
	return false
}

func isLuckyStroke(stroke int) bool {
	luckyStrokes := []int{1, 3, 5, 6, 7, 8, 11, 13, 15, 16, 18, 21, 23, 24, 25, 29, 31}
	for _, s := range luckyStrokes {
		if s == stroke {
			return true
		}
	}
	return false
}

func getTone(pinyin string) int {
	if len(pinyin) == 0 {
		return 0
	}
	last := pinyin[len(pinyin)-1]
	if last >= '1' && last <= '4' {
		return int(last - '0')
	}
	return 0
}

func getShengMu(pinyin string) string {
	if len(pinyin) == 0 {
		return ""
	}
	shengmuList := []string{"zh", "ch", "sh", "b", "p", "m", "f", "d", "t", "n", "l", "g", "k", "h", "j", "q", "x", "z", "c", "s", "r", "y", "w"}
	for _, sm := range shengmuList {
		if len(pinyin) >= len(sm) && pinyin[:len(sm)] == sm {
			return sm
		}
	}
	return ""
}

func getYunMu(pinyin string) string {
	sm := getShengMu(pinyin)
	if sm == "" {
		return pinyin
	}
	return pinyin[len(sm):]
}

func joinNotes(notes []string) string {
	if len(notes) == 0 {
		return ""
	}
	result := notes[0]
	for i := 1; i < len(notes); i++ {
		result += "; " + notes[i]
	}
	return result
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func firstPinyin(py []string) string {
	if len(py) > 0 {
		return py[0]
	}
	return ""
}

type PoetryRater struct {
	cfg        *config.Config
	poetryChar map[string]bool
}

func NewPoetryRater(cfg *config.Config, poetryChars map[string]bool) *PoetryRater {
	return &PoetryRater{cfg: cfg, poetryChar: poetryChars}
}

func (r *PoetryRater) Name() string {
	return "poetry"
}

func (r *PoetryRater) Weight() float64 {
	return r.cfg.Rate.PoetryWeight
}

func (r *PoetryRater) Rate(name *NameInfo, _ *chronosfate.FateData) (float64, string) {
	if len(r.poetryChar) == 0 {
		return 70.0, "诗词数据未加载"
	}

	score := 50.0
	var notes []string

	hasPoetry1 := r.poetryChar[name.Char1.Char]
	hasPoetry2 := r.poetryChar[name.Char2.Char]

	if hasPoetry1 {
		score += 20.0
		notes = append(notes, name.Char1.Char+"出自诗词")
	}
	if hasPoetry2 {
		score += 20.0
		notes = append(notes, name.Char2.Char+"出自诗词")
	}

	if hasPoetry1 && hasPoetry2 {
		score += 10.0
		notes = append(notes, "双字皆有诗词出处")
	}

	if !hasPoetry1 && !hasPoetry2 {
		score -= 10.0
		notes = append(notes, "无诗词出处")
	}

	if score > 100.0 {
		score = 100.0
	}
	if score < 0.0 {
		score = 0.0
	}

	var noteStr string
	if len(notes) > 0 {
		noteStr = joinNotes(notes)
	} else {
		noteStr = "诗词中和"
	}

	return score, noteStr
}
