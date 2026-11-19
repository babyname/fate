package naming

import (
	"strings"
	"testing"

	"github.com/babyname/fate/v4/ent"
	"github.com/babyname/fate/v4/internal/analysis"
)

func TestNaturalNameSoundIncludesSurname(t *testing.T) {
	for _, tc := range []struct {
		readings   []string
		acceptable bool
		notice     string
	}{
		{[]string{"máo", "míng", "qiān"}, true, ""},
		{[]string{"máo", "bǐng", "yán"}, true, "毛病"},
		{[]string{"máo", "bìng", "yán"}, false, ""},
		{[]string{"dù", "zi", "téng"}, false, ""},
		{[]string{"lǐ", "lǐ", "níng"}, true, "变调"},
		{[]string{"zhāng", "shū", "níng"}, true, ""},
	} {
		parts, ok := scoreNaturalName(tc.readings, &ent.Character{Char: "明"}, &ent.Character{Char: "谦"}, true)
		if ok != tc.acceptable {
			t.Errorf("scoreNaturalName(%v) accepted = %v, want %v", tc.readings, ok, tc.acceptable)
		}
		if ok && tc.notice != "" && !strings.Contains(strings.Join(append(parts[0].Reasons, parts[1].Reasons...), " "), tc.notice) {
			t.Errorf("scoreNaturalName(%v) missing notice %q", tc.readings, tc.notice)
		}
	}
	if pinyinInitial("zhang") == pinyinInitial("zi") || parseNameSound("lü4") != parseNameSound("lǜ") {
		t.Error("pinyin normalization confused initials or tone spellings")
	}
}

func TestNaturalNameMeaningAndUsabilityUseCandidateCharacters(t *testing.T) {
	first := &ent.Character{Char: "安", Regular: true, Pinyin: []string{"ān"}}
	compatible := &ent.Character{Char: "宁", Regular: true, Pinyin: []string{"níng", "nìng"}}
	unrelated := &ent.Character{Char: "山", Regular: true, Pinyin: []string{"shān"}}
	if scoreMeaning(first, compatible, false).Score <= scoreMeaning(first, unrelated, false).Score {
		t.Error("a shared theme should affect meaning even without a curated description")
	}
	if scoreNameUsability(first, compatible).Score >= scoreNameUsability(first, unrelated).Score {
		t.Error("an alternative pronunciation should be surfaced as a reading risk")
	}
}

func TestCombinedRankingUsesTraditionalDimensions(t *testing.T) {
	sound := []analysis.RecommendationScorePart{{Score: 86, Weight: 1}}
	traditional := analysis.ScoreDetail{WuGeShuLi: 70, WuXingBaZi: 98, ShengXiao: 94}
	withoutBirth := combinedRanking(sound, traditional, false)
	if withoutBirth.Traditional != 70 || withoutBirth.Total != 79.6 || withoutBirth.UsesBirth {
		t.Fatalf("five-grid-only ranking = %+v", withoutBirth)
	}
	withBirth := combinedRanking(sound, traditional, true)
	if withBirth.Traditional != 84.5 || withBirth.Total != 85.4 || !withBirth.UsesBirth {
		t.Fatalf("ranking with birth data = %+v", withBirth)
	}
	higherTraditional := combinedRanking(sound, analysis.ScoreDetail{WuGeShuLi: 90}, false)
	if higherTraditional.Total <= withoutBirth.Total {
		t.Fatalf("equal sound and meaning must rank higher with a higher traditional score")
	}
}
