package seeddb

import (
	"fmt"

	"github.com/babyname/fate/internal/wuge"
)

func transformCharacters(oldList []oldCharacter) []SeedCharacter {
	result := make([]SeedCharacter, 0, len(oldList))
	for _, old := range oldList {
		result = append(result, transformCharacter(old))
	}
	return result
}

func transformCharacter(old oldCharacter) SeedCharacter {
	sc := SeedCharacter{
		Char:      old.Ch,
		IsKangxi:  old.IsKangXi,
		Regular:   old.Regular,
		Nameable:  old.NameScience,
		WuXing:    old.WuXing,
		Source:    "legacy",
	}

	sc.Pinyin = parseJSONString(old.PinYin)

	if old.Radical != "" {
		sc.Radical = old.Radical
		sc.RadicalStroke = old.RadicalStroke
	}

	sc.SimplifiedStroke = old.SimpleTotalStroke
	sc.TraditionalStroke = old.TraditionalTotalStroke
	sc.KangxiStroke = old.KangXiStroke
	sc.ScienceStroke = old.ScienceStroke

	tradChars := parseJSONString(old.TraditionalCharacter)
	if len(tradChars) > 0 {
		sc.IsSimplified = true
		sc.SimplifiedOfChar = tradChars[0]
	}

	if old.TraditionalTotalStroke > 0 && old.SimpleTotalStroke != old.TraditionalTotalStroke && len(tradChars) == 0 {
		sc.IsTraditional = true
	}

	if old.KangXi != "" && old.KangXi != old.Ch {
		sc.IsVariant = true
		sc.VariantOfChar = old.KangXi
	}

	variantChars := parseJSONString(old.VariantCharacter)
	if len(variantChars) > 0 {
		sc.IsVariant = true
		if sc.VariantOfChar == "" {
			sc.VariantOfChar = variantChars[0]
		}
	}

	commentParts := parseJSONString(old.Comment)
	if old.Lucky != "" {
		sc.Comment = fmt.Sprintf("lucky=%s", old.Lucky)
	}
	if len(commentParts) > 0 {
		if sc.Comment != "" {
			sc.Comment += "; "
		}
		if len(commentParts) == 1 {
			if len(commentParts[0]) > 200 {
				sc.Comment += commentParts[0][:200] + "..."
			} else {
				sc.Comment += commentParts[0]
			}
		}
	}

	return sc
}

func transformWuGeLucky(oldList []oldWuGeLucky) []SeedWuGeLucky {
	result := make([]SeedWuGeLucky, 0, len(oldList))
	for _, old := range oldList {
		result = append(result, transformWuGeLuckyOne(old))
	}
	return result
}

func transformWuGeLuckyOne(old oldWuGeLucky) SeedWuGeLucky {
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

func transformWuXing(oldList []oldWuXing) []SeedWuXing {
	result := make([]SeedWuXing, 0, len(oldList))
	for _, old := range oldList {
		result = append(result, SeedWuXing{
			ID:      old.ID,
			First:   old.SanCai,
			Second:  old.SanCai,
			Third:   old.SanCai,
			Fortune: formatWuXingFortune(old.Lucky, old.Comment),
		})
	}
	return result
}

func formatWuXingFortune(lucky bool, comment string) string {
	if lucky {
		return "吉|" + comment
	}
	return "凶|" + comment
}
