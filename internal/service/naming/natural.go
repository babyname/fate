package naming

import (
	"context"
	_ "embed"
	"math"
	"slices"
	"sort"
	"strings"

	"github.com/babyname/fate/v4/ent"
	"github.com/babyname/fate/v4/ent/character"
	"github.com/babyname/fate/v4/internal/analysis"
	"github.com/babyname/fate/v4/internal/chronosfate"
	"github.com/babyname/fate/v4/internal/naming"
)

//go:embed natural_names.tsv
var naturalNameData string

type naturalName struct {
	first, second string
	preference    string
	meaning       string
}

var naturalNames = loadNaturalNames()

func loadNaturalNames() []naturalName {
	var names []naturalName
	for _, line := range strings.Split(naturalNameData, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Split(line, "\t")
		if len(parts) != 3 {
			panic("invalid curated name row: " + line)
		}
		pair := []rune(parts[1])
		if len(pair) != 2 || parts[2] == "" || (parts[0] != "男" && parts[0] != "女" && parts[0] != "中") {
			panic("invalid curated name row: " + line)
		}
		names = append(names, naturalName{string(pair[0]), string(pair[1]), parts[0], parts[2]})
	}
	return names
}

// GenerateRecommendations scores the curated catalogue with the given surname.
func (p *Pipeline) GenerateRecommendations(ctx context.Context, req GenerateRequest, required, avoided []string) ([]analysis.NameResult, *ExcellentTable, error) {
	preference := "男"
	if req.Sex == naming.SexGirl {
		preference = "女"
	}
	surname := surnameStr(req.LastName)
	seeds := appendRequiredCombinations(slices.Clone(naturalNames), required)
	eligible := make([]naturalName, 0, len(seeds))
	requested := make(map[string]bool)
	for _, name := range seeds {
		if name.preference != "中" && name.preference != preference {
			continue
		}
		if !allowsNaturalPair(name, surname, required, avoided) {
			continue
		}
		eligible = append(eligible, name)
		requested[name.first] = true
		requested[name.second] = true
	}

	lookup := make([]string, 0, len(requested))
	for char := range requested {
		lookup = append(lookup, char)
	}
	chars, err := p.db.Character.Query().Where(character.CharIn(lookup...)).All(ctx)
	if err != nil {
		return nil, nil, err
	}
	charMap := make(map[string]*ent.Character, len(chars))
	for _, char := range chars {
		if reading := givenNameReadings[char.Char]; reading != "" {
			copy := *char
			copy.Pinyin = append([]string{reading}, char.Pinyin...)
			char = &copy
		}
		charMap[char.Char] = char
	}

	var fateData *chronosfate.FateData
	if !req.Born.IsZero() {
		fateData = p.resolveFateData(req)
	}
	strokes := surnameStrokes(req.Filter, req.LastName)
	surnameSounds := surnamePinyin(req.LastName)
	candidates := make([]analysis.NameResult, 0, len(eligible))
	table := NewExcellentTable()
	seen := make(map[string]bool)
	for _, name := range eligible {
		key := name.first + name.second
		if seen[key] {
			continue
		}
		c1, c2 := charMap[name.first], charMap[name.second]
		p1, p2 := firstPinyin(c1), firstPinyin(c2)
		if p1 == "" || p2 == "" || c1.ScienceStroke <= 0 || c2.ScienceStroke <= 0 {
			continue
		}
		if req.Filter.PoetryMode() == 2 && (!c1.HasPoetry || !c2.HasPoetry) {
			continue
		}
		readings := append(slices.Clone(surnameSounds), p1, p2)
		sound, acceptable := scoreNaturalName(readings, c1, c2, name.meaning != requiredMeaning)
		if !acceptable {
			continue
		}
		seen[key] = true
		result := analysis.BuildNameResult(0, surname, c1, c2, strokes[0], strokes[1], fateData)
		ranking := combinedRanking(sound, result.ScoreDetail, fateData != nil)
		result.Ranking = &ranking
		result.Recommendation = name.meaning
		result.FullPinyin = strings.Join(readings, " ")
		candidates = append(candidates, result)
		table.TryPush(ExcellentEntry{
			Char1: c1.Char, Char2: c2.Char, Score: ranking.Total, WuXing1: c1.WuXing, WuXing2: c2.WuXing,
			HasPoetry: result.HasPoetry,
		})
	}
	sort.SliceStable(candidates, func(i, j int) bool { return candidates[i].Ranking.Total > candidates[j].Ranking.Total })
	for i := range candidates {
		candidates[i].Rank = i + 1
	}
	table.Finalize()
	return candidates, table, nil
}

func allowsNaturalPair(name naturalName, surname string, required, avoided []string) bool {
	if name.first == name.second || strings.Contains(surname, name.first) || strings.Contains(surname, name.second) {
		return false
	}
	if slices.Contains(avoided, name.first) || slices.Contains(avoided, name.second) {
		return false
	}
	for _, char := range required {
		if char != name.first && char != name.second {
			return false
		}
	}
	return true
}

const requiredMeaning = "按指定字搭配，可结合字义和读音挑选。"

func appendRequiredCombinations(names []naturalName, required []string) []naturalName {
	switch len(required) {
	case 1:
		for _, char := range "安宁远和言舟清初明川书然" {
			names = append(names,
				naturalName{required[0], string(char), "中", requiredMeaning},
				naturalName{string(char), required[0], "中", requiredMeaning},
			)
		}
	case 2:
		names = append(names,
			naturalName{required[0], required[1], "中", requiredMeaning},
			naturalName{required[1], required[0], "中", requiredMeaning},
		)
	}
	return names
}

func combinedRanking(sound []analysis.RecommendationScorePart, detail analysis.ScoreDetail, usesBirth bool) analysis.RecommendationRanking {
	traditionalParts := []analysis.RecommendationScorePart{
		{Label: "五格与三才", Score: detail.WuGeShuLi, Weight: 1},
	}
	if usesBirth {
		traditionalParts[0].Weight = 30.0 / 65
		traditionalParts = append(traditionalParts,
			analysis.RecommendationScorePart{Label: "五行八字", Score: detail.WuXingBaZi, Weight: 25.0 / 65},
			analysis.RecommendationScorePart{Label: "生肖", Score: detail.ShengXiao, Weight: 10.0 / 65},
		)
	}
	soundScore := weightedScore(sound)
	traditionalScore := weightedScore(traditionalParts)
	return analysis.RecommendationRanking{
		SoundAndMeaning: soundScore, Traditional: traditionalScore,
		Total: math.Round((0.6*soundScore+0.4*traditionalScore)*10) / 10, UsesBirth: usesBirth,
		SoundDetail: sound, TraditionalDetail: traditionalParts,
	}
}

func weightedScore(parts []analysis.RecommendationScorePart) float64 {
	var total float64
	for _, part := range parts {
		total += part.Score * part.Weight
	}
	return math.Round(total*10) / 10
}

var surnameReadings = map[string]string{
	"张": "zhāng", "張": "zhāng", "王": "wáng", "李": "lǐ", "毛": "máo",
	"单": "shàn", "單": "shàn", "解": "xiè", "曾": "zēng", "仇": "qiú",
	"查": "zhā", "朴": "piáo", "区": "ōu", "區": "ōu",
	"欧阳": "ōu yáng", "歐陽": "ōu yáng", "司马": "sī mǎ", "司馬": "sī mǎ",
	"诸葛": "zhū gě", "諸葛": "zhū gě", "尉迟": "yù chí", "尉遲": "yù chí",
	"长孙": "zhǎng sūn", "長孫": "zhǎng sūn", "万俟": "mò qí", "萬俟": "mò qí",
}

var givenNameReadings = map[string]string{
	"宁": "níng", "予": "yǔ", "和": "hé", "清": "qīng", "澄": "chéng",
	"荷": "hé", "栖": "qī", "杉": "shān", "蹊": "xī", "冠": "guān",
}

func surnamePinyin(chars [2]*ent.Character) []string {
	if reading := surnameReadings[surnameStr(chars)]; reading != "" {
		return strings.Fields(reading)
	}
	var readings []string
	for _, char := range chars {
		if reading := firstPinyin(char); reading != "" {
			readings = append(readings, reading)
		}
	}
	return readings
}

func firstPinyin(char *ent.Character) string {
	if char == nil || len(char.Pinyin) == 0 {
		return ""
	}
	return strings.TrimSpace(char.Pinyin[0])
}
