package main

import (
	"fmt"
	"strings"

	"github.com/babyname/fate/internal/wuge"
)

type SeedCharacter struct {
	ID                 int      `json:"id"`
	Char               string   `json:"char"`
	Unicode            string   `json:"unicode,omitempty"`
	IsSimplified       bool     `json:"is_simplified"`
	IsTraditional      bool     `json:"is_traditional"`
	IsKangxi           bool     `json:"is_kangxi"`
	IsVariant          bool     `json:"is_variant"`
	IsAncient          bool     `json:"is_ancient"`
	Pinyin             []string `json:"pinyin,omitempty"`
	Radical            string   `json:"radical,omitempty"`
	RadicalStroke      int      `json:"radical_stroke,omitempty"`
	SimplifiedStroke   int      `json:"simplified_stroke,omitempty"`
	TraditionalStroke  int      `json:"traditional_stroke,omitempty"`
	KangxiStroke       int      `json:"kangxi_stroke,omitempty"`
	ScienceStroke      int      `json:"science_stroke,omitempty"`
	WuXing             string   `json:"wu_xing,omitempty"`
	Regular            bool     `json:"regular"`
	CommonLevel        int      `json:"common_level,omitempty"`
	GenderHint         string   `json:"gender_hint,omitempty"`
	Nameable           bool     `json:"nameable"`
	Meaning            string   `json:"meaning,omitempty"`
	Source             string   `json:"source,omitempty"`
	SourceConfidence   float64  `json:"source_confidence,omitempty"`
	Comment            string   `json:"comment,omitempty"`
	SimplifiedOfChar   string   `json:"simplified_of_char,omitempty"`
	TraditionalOfChar  string   `json:"traditional_of_char,omitempty"`
	VariantOfChar      string   `json:"variant_of_char,omitempty"`
}

type SeedWuGeLucky struct {
	LastStroke1  int    `json:"last_stroke_1"`
	LastStroke2  int    `json:"last_stroke_2"`
	FirstStroke1 int    `json:"first_stroke_1"`
	FirstStroke2 int    `json:"first_stroke_2"`
	TianGe       int    `json:"tian_ge"`
	TianDaYan    string `json:"tian_da_yan"`
	RenGe        int    `json:"ren_ge"`
	RenDaYan     string `json:"ren_da_yan"`
	DiGe         int    `json:"di_ge"`
	DiDaYan      string `json:"di_da_yan"`
	WaiGe        int    `json:"wai_ge"`
	WaiDaYan     string `json:"wai_da_yan"`
	ZongGe       int    `json:"zong_ge"`
	ZongDaYan    string `json:"zong_da_yan"`
	ZongLucky    bool   `json:"zong_lucky"`
	ZongSex      bool   `json:"zong_sex"`
	ZongMax      bool   `json:"zong_max"`
}

type SeedWuXing struct {
	ID      string `json:"id"`
	First   string `json:"first"`
	Second  string `json:"second"`
	Third   string `json:"third"`
	Fortune string `json:"fortune"`
}

func TransformCharacters(oldList []interface{}) []SeedCharacter {
	result := make([]SeedCharacter, 0, len(oldList))
	for _, item := range oldList {
		old := item.(OldCharacter)
		sc := transformCharacter(old)
		result = append(result, sc)
	}
	return result
}

func transformCharacter(old OldCharacter) SeedCharacter {
	sc := SeedCharacter{
		ID:        old.ID,
		Char:      old.Ch,
		IsKangxi:  old.IsKangXi,
		Regular:   old.Regular,
		Nameable:  old.NameScience,
		WuXing:    old.WuXing,
		Source:    "legacy",
	}

	sc.Pinyin = parsePinyin(old.PinYin)

	sc.SimplifiedStroke = old.SimpleTotalStroke
	sc.TraditionalStroke = old.TraditionalTotalStroke
	sc.KangxiStroke = old.KangXiStroke
	sc.ScienceStroke = old.ScienceStroke

	if old.TraditionalCharacter != "" && old.TraditionalCharacter != old.Ch {
		sc.IsSimplified = true
		sc.SimplifiedOfChar = old.TraditionalCharacter
	}
	if old.TraditionalCharacter == "" && old.SimpleTotalStroke != old.TraditionalTotalStroke && old.TraditionalTotalStroke > 0 {
		sc.IsTraditional = true
	}
	if old.KangXi != "" && old.KangXi != old.Ch {
		sc.IsVariant = true
		sc.VariantOfChar = old.KangXi
	}
	if old.VariantCharacter != "" && old.VariantCharacter != old.Ch {
		sc.IsVariant = true
		if sc.VariantOfChar == "" {
			sc.VariantOfChar = old.VariantCharacter
		}
	}

	if old.Lucky != "" {
		sc.Comment = fmt.Sprintf("lucky=%s", old.Lucky)
	}

	return sc
}

func parsePinyin(pinyin string) []string {
	if pinyin == "" {
		return nil
	}
	parts := strings.Split(pinyin, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func TransformWuGeLucky(oldList []interface{}) []SeedWuGeLucky {
	result := make([]SeedWuGeLucky, 0, len(oldList))
	for _, item := range oldList {
		old := item.(OldWuGeLucky)
		sw := transformWuGeLucky(old)
		result = append(result, sw)
	}
	return result
}

func transformWuGeLucky(old OldWuGeLucky) SeedWuGeLucky {
	ge := wuge.CalcWuGe(old.LastStroke1, old.LastStroke2, old.FirstStroke1, old.FirstStroke2)

	tianDaYan := wuge.Find(ge.TianGe())
	renDaYan := wuge.Find(ge.RenGe())
	diDaYan := wuge.Find(ge.DiGe())
	waiDaYan := wuge.Find(ge.WaiGe())
	zongDaYan := wuge.Find(ge.ZongGe())

	return SeedWuGeLucky{
		LastStroke1:  old.LastStroke1,
		LastStroke2:  old.LastStroke2,
		FirstStroke1: old.FirstStroke1,
		FirstStroke2: old.FirstStroke2,
		TianGe:       ge.TianGe(),
		TianDaYan:    formatDaYan(tianDaYan),
		RenGe:        ge.RenGe(),
		RenDaYan:     formatDaYan(renDaYan),
		DiGe:         ge.DiGe(),
		DiDaYan:      formatDaYan(diDaYan),
		WaiGe:        ge.WaiGe(),
		WaiDaYan:     formatDaYan(waiDaYan),
		ZongGe:       ge.ZongGe(),
		ZongDaYan:    formatDaYan(zongDaYan),
		ZongLucky:    old.ZongLucky,
		ZongSex:      bool(zongDaYan.Sex),
		ZongMax:      zongDaYan.Max,
	}
}

func formatDaYan(dy wuge.DaYan) string {
	return fmt.Sprintf("%d|%s|%s", dy.Number, dy.Lucky, dy.SkyNine)
}

func TransformWuXing(oldList []interface{}) []SeedWuXing {
	result := make([]SeedWuXing, 0, len(oldList))
	for _, item := range oldList {
		old := item.(OldWuXing)
		sw := transformWuXing(old)
		result = append(result, sw)
	}
	return result
}

func transformWuXing(old OldWuXing) SeedWuXing {
	return SeedWuXing{
		ID:      old.ID,
		First:   old.SanCai,
		Second:  old.SanCai,
		Third:   old.SanCai,
		Fortune: formatWuXingFortune(old.Lucky, old.Comment),
	}
}

func formatWuXingFortune(lucky bool, comment string) string {
	if lucky {
		return "吉|" + comment
	}
	return "凶|" + comment
}
