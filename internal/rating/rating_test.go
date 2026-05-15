package rating

import (
	"testing"
	"time"

	"github.com/babyname/fate/ent"
	v2 "github.com/godcong/chronos/v2"
)

func TestNewRater(t *testing.T) {
	fateData, _ := v2.GetFateData(&v2.FateInput{
		BirthDate: time.Date(2024, 1, 15, 10, 30, 0, 0, time.Local),
		Gender:    1,
		Surname:   "张",
	})

	rater := NewRater(fateData)
	if rater == nil {
		t.Fatal("NewRater should not return nil")
	}
}

func TestRateName(t *testing.T) {
	fateData, _ := v2.GetFateData(&v2.FateInput{
		BirthDate: time.Date(2024, 1, 15, 10, 30, 0, 0, time.Local),
		Gender:    1,
		Surname:   "张",
	})

	rater := NewRaterWithStrokes(fateData, 11, 0)

	c1 := &ent.Character{
		Char:          "伟",
		WuXing:        "土",
		ScienceStroke: 11,
		Pinyin:        []string{"wěi"},
	}
	c2 := &ent.Character{
		Char:          "安",
		WuXing:        "土",
		ScienceStroke: 6,
		Pinyin:        []string{"ān"},
	}

	rating := rater.RateName("张", c1, c2)
	if rating == nil {
		t.Fatal("RateName should not return nil")
	}

	if rating.TotalScore <= 0 || rating.TotalScore > 100 {
		t.Errorf("TotalScore should be between 0 and 100, got %.1f", rating.TotalScore)
	}
	if rating.WuXingScore <= 0 || rating.WuXingScore > 100 {
		t.Errorf("WuXingScore should be between 0 and 100, got %.1f", rating.WuXingScore)
	}
	if rating.WuGeScore <= 0 || rating.WuGeScore > 100 {
		t.Errorf("WuGeScore should be between 0 and 100, got %.1f", rating.WuGeScore)
	}
	if rating.YinYunScore <= 0 || rating.YinYunScore > 100 {
		t.Errorf("YinYunScore should be between 0 and 100, got %.1f", rating.YinYunScore)
	}
	if rating.Grade == "" {
		t.Error("Grade should not be empty")
	}
	if rating.Interpret == "" {
		t.Error("Interpret should not be empty")
	}

	t.Logf("Total: %.1f (%s)", rating.TotalScore, rating.Grade)
	t.Logf("WenHua: %.1f - %s", rating.WenHuaScore, rating.WenHuaDetail)
	t.Logf("WuXing: %.1f - %s", rating.WuXingScore, rating.WuXingDetail)
	t.Logf("ShengXiao: %.1f - %s", rating.ShengXiaoScore, rating.ShengXiaoDetail)
	t.Logf("WuGe: %.1f - %s", rating.WuGeScore, rating.WuGeDetail)
	t.Logf("YinYun: %.1f - %s", rating.YinYunScore, rating.YinYunDetail)
	t.Logf("Interpret: %s", rating.Interpret)
}

func TestScoreToGrade(t *testing.T) {
	tests := []struct {
		score float64
		grade string
	}{
		{95, "上上"},
		{85, "上吉"},
		{75, "中吉"},
		{65, "中平"},
		{55, "中下"},
		{40, "下下"},
	}

	for _, tt := range tests {
		got := scoreToGrade(tt.score)
		if got != tt.grade {
			t.Errorf("scoreToGrade(%.0f) = %s, want %s", tt.score, got, tt.grade)
		}
	}
}

func TestRateNameWithXiWuxing(t *testing.T) {
	fateData, _ := v2.GetFateData(&v2.FateInput{
		BirthDate: time.Date(2024, 1, 15, 10, 30, 0, 0, time.Local),
		Gender:    1,
		Surname:   "张",
	})

	rater := NewRaterWithStrokes(fateData, 11, 0)

	c1 := &ent.Character{
		Char:          "浩",
		WuXing:        "水",
		ScienceStroke: 11,
		Pinyin:        []string{"hào"},
	}
	c2 := &ent.Character{
		Char:          "然",
		WuXing:        "金",
		ScienceStroke: 12,
		Pinyin:        []string{"rán"},
	}

	rating := rater.RateName("张", c1, c2)
	t.Logf("浩然: Total=%.1f, WuXing=%.1f, WuGe=%.1f, Grade=%s", rating.TotalScore, rating.WuXingScore, rating.WuGeScore, rating.Grade)
	t.Logf("Interpret: %s", rating.Interpret)
}
