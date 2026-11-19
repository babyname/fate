package naming

import (
	"fmt"
	"slices"
	"strings"
	"unicode"

	"github.com/babyname/fate/v4/ent"
	"github.com/babyname/fate/v4/internal/analysis"
	"golang.org/x/text/unicode/norm"
)

type nameSound struct {
	base string
	tone int
}

func parseNameSound(reading string) nameSound {
	var base strings.Builder
	var tone int
	for _, letter := range norm.NFD.String(strings.ToLower(reading)) {
		switch letter {
		case '\u0304':
			tone = 1
		case '\u0301':
			tone = 2
		case '\u030c':
			tone = 3
		case '\u0300':
			tone = 4
		case '1', '2', '3', '4':
			tone = int(letter - '0')
		case '\u0308':
			base.WriteRune(letter)
		case 'v':
			base.WriteRune('ü')
		default:
			if unicode.IsLetter(letter) {
				base.WriteRune(letter)
			}
		}
	}
	return nameSound{norm.NFC.String(base.String()), tone}
}

func pinyinInitial(base string) string {
	for _, initial := range []string{"zh", "ch", "sh"} {
		if strings.HasPrefix(base, initial) {
			return initial
		}
	}
	if base != "" && strings.ContainsRune("bpmfdtnlgkhjqxzcsrwy", rune(base[0])) {
		return base[:1]
	}
	return ""
}

func scoreNaturalName(readings []string, first, second *ent.Character, hasPairMeaning bool) ([]analysis.RecommendationScorePart, bool) {
	sounds := make([]nameSound, len(readings))
	for i, reading := range readings {
		sounds[i] = parseNameSound(reading)
		if sounds[i].base == "" {
			return nil, false
		}
	}
	homophones, acceptable := scoreHomophones(sounds)
	if !acceptable {
		return nil, false
	}
	return []analysis.RecommendationScorePart{
		scorePronunciation(sounds),
		homophones,
		scoreMeaning(first, second, hasPairMeaning),
		scoreNameUsability(first, second),
	}, true
}

func scorePronunciation(sounds []nameSound) analysis.RecommendationScorePart {
	part := analysis.RecommendationScorePart{Label: "全名连读", Score: 94, Weight: .35}
	var repeatedSyllable, repeatedInitial, repeatedFinal, repeatedTone, thirdTones bool
	for i := 1; i < len(sounds); i++ {
		previous, next := sounds[i-1], sounds[i]
		if previous.base == next.base {
			part.Score -= 9
			repeatedSyllable = true
		} else {
			firstInitial, secondInitial := pinyinInitial(previous.base), pinyinInitial(next.base)
			if firstInitial != "" && firstInitial == secondInitial {
				part.Score -= 3
				repeatedInitial = true
			}
			if previous.base[len(firstInitial):] == next.base[len(secondInitial):] {
				part.Score -= 3
				repeatedFinal = true
			}
		}
		if previous.tone == 3 && next.tone == 3 {
			thirdTones = true
		} else if previous.tone != 0 && previous.tone == next.tone {
			part.Score -= 2
			repeatedTone = true
		}
	}
	if repeatedSyllable {
		part.Reasons = append(part.Reasons, "相邻两字音节相同，读全名确认是否容易混听。")
	}
	if repeatedInitial {
		part.Reasons = append(part.Reasons, "相邻字有相同声母。")
	}
	if repeatedFinal {
		part.Reasons = append(part.Reasons, "相邻字韵母相同，有押韵感。")
	}
	if repeatedTone {
		part.Reasons = append(part.Reasons, "相邻字有同调，语调变化较少。")
	}
	if thirdTones {
		part.Reasons = append(part.Reasons, "相邻三声连读时常会变调，建议按全名试读。")
	}
	if len(part.Reasons) == 0 {
		part.Reasons = []string{"相邻字的音节和声调有变化。"}
	}
	return part
}

type nameSoundRisk struct {
	word  string
	tones string
}

var riskyNameSounds = map[string]nameSoundRisk{
	"mao bing": {"毛病", "24"}, "zhi zhang": {"智障", "44"}, "yi han": {"遗憾", "24"},
	"shi bai": {"失败", "14"}, "wu de": {"无德", "22"}, "wu li": {"无理", "23"},
	"wu nai": {"无奈", "24"}, "wang ba": {"王八", "21"}, "fan tong": {"饭桶", "43"},
	"du zi teng": {"肚子疼", "402"}, "shi zhen xiang": {"屎真香", "311"},
	"zhu yi qun": {"猪一群", "142"}, "hu li jing": {"狐狸精", "201"}, "yang wei": {"阳痿", "23"},
}

func scoreHomophones(sounds []nameSound) (analysis.RecommendationScorePart, bool) {
	part := analysis.RecommendationScorePart{Label: "谐音风险", Score: 96, Weight: .20}
	bases := make([]string, len(sounds))
	for i, sound := range sounds {
		bases[i] = sound.base
	}
	for start := range sounds {
		for end := start + 2; end <= min(len(sounds), start+3); end++ {
			risk, found := riskyNameSounds[strings.Join(bases[start:end], " ")]
			if !found {
				continue
			}
			mismatches := 0
			for i, sound := range sounds[start:end] {
				if byte('0'+sound.tone) != risk.tones[i] {
					mismatches++
				}
			}
			if mismatches == 0 {
				return part, false
			}
			score := 54.0
			if mismatches > 1 {
				score = 77
			}
			if score < part.Score {
				part.Score = score
				part.Reasons = []string{fmt.Sprintf("与「%s」音节相近、声调有别，建议全名试读。", risk.word)}
			}
		}
	}
	if len(part.Reasons) == 0 {
		part.Reasons = []string{"未命中普通话常见负面谐音词表。"}
	}
	return part, true
}

var meaningThemes = []struct{ label, chars string }{
	{"安定从容", "安宁和舒静悠然庭澄怡平晏宜宛婉令清慎"},
	{"求知修养", "书文明知思学修谦序衡睿正雅砚诗微见问理则仪"},
	{"品格心意", "诚信真允仁德善慈敬毅谦嘉念立恒意心诺温恺"},
	{"光景时节", "明景朗晴晓晗昕星月晨时辰初夏秋冬春晖光照映曦"},
	{"山水草木", "山岳岚川江溪泽海水舟帆洲屿云雪雨松杉荷兰棠桐芷薇树禾波澜"},
	{"志向与成长", "远遥启初致承鸿凌志行卓弘嘉新高长腾鹏"},
}

func scoreMeaning(first, second *ent.Character, hasPairMeaning bool) analysis.RecommendationScorePart {
	part := analysis.RecommendationScorePart{Label: "字义搭配", Score: 62, Weight: .30}
	var themes [2][]string
	for i, char := range []*ent.Character{first, second} {
		for _, theme := range meaningThemes {
			if strings.Contains(theme.chars, char.Char) {
				themes[i] = append(themes[i], theme.label)
			}
		}
		if len(themes[i]) > 0 {
			part.Score += 8
		} else if char.Meaning != "" {
			part.Score += 3
		}
		if strings.Contains("病死恨臭丑凶虐灾怨毒屎厄", char.Char) {
			part.Score -= 24
			part.Reasons = append(part.Reasons, fmt.Sprintf("「%s」的常见含义容易带来负面联想。", char.Char))
		}
	}
	for _, theme := range themes[0] {
		if slices.Contains(themes[1], theme) {
			part.Score += 10
			part.Reasons = append(part.Reasons, fmt.Sprintf("两字可共同表达「%s」。", theme))
			break
		}
	}
	if hasPairMeaning {
		part.Score += 6
		part.Reasons = append(part.Reasons, "可结合候选的取名寓意判断两字搭配。")
	}
	if len(part.Reasons) == 0 {
		part.Reasons = []string{"建议结合详情里的单字释义确认两字能否表达同一个心愿。"}
	}
	part.Score = max(0, part.Score)
	return part
}

func scoreNameUsability(first, second *ent.Character) analysis.RecommendationScorePart {
	part := analysis.RecommendationScorePart{Label: "易读易写", Score: 96, Weight: .15}
	for _, char := range []*ent.Character{first, second} {
		readings := make(map[nameSound]bool)
		for _, reading := range char.Pinyin {
			sound := parseNameSound(reading)
			if sound.base != "" {
				readings[sound] = true
			}
		}
		if len(readings) > 1 {
			part.Score -= 5
			part.Reasons = append(part.Reasons, fmt.Sprintf("「%s」有其他读音，此处取 %s。", char.Char, firstPinyin(char)))
		}
		if char.CommonLevel >= 3 || !char.Regular {
			part.Score -= 6
			part.Reasons = append(part.Reasons, fmt.Sprintf("「%s」较少见，日常可能需要解释用字。", char.Char))
		}
		stroke := char.SimplifiedStroke
		if !char.IsSimplified || stroke <= 0 {
			stroke = char.TraditionalStroke
		}
		if stroke > 14 {
			part.Score -= float64(min(10, stroke-14))
			part.Reasons = append(part.Reasons, fmt.Sprintf("「%s」书写笔画较多（%d 画）。", char.Char, stroke))
		}
	}
	if len(part.Reasons) == 0 {
		part.Reasons = []string{"两字的常用度、读音和书写笔画适合日常使用。"}
	}
	return part
}
