package fate

import (
	"testing"
	"time"

	v2 "github.com/godcong/chronos/v2"
	"github.com/babyname/fate/config"
	"github.com/babyname/fate/ent"
	"github.com/babyname/fate/internal/analysis"
)

func TestNameAnalysisOutput(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Database = config.DBConfig{
		Driver: "sqlite3",
		Name:   testDBPath,
	}

	f, err := New(cfg)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	born, _ := time.Parse("2006/01/02 15:04", "2024/06/15 10:30")

	s := f.NewSessionWithFilter(NewFilter(FilterOption{
		CharacterFilter:     true,
		CharacterFilterType: CharacterFilterTypeDefault,
		MinStroke:           3,
		MaxStroke:           18,
		RegularFilter:       true,
		DaYanFilter:         true,
		WuXingFilter:        true,
	}))

	input := &Input{
		Last: [2]string{"张", ""},
		Born: born,
		Sex:  Sex(1),
	}

	start := time.Now()
	err = s.Start(input)
	if err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	s.Wait()
	elapsed := time.Since(start)

	output := input.Output()
	total := output.Total()
	t.Logf("共生成 %d 个名字, 耗时 %v", total, elapsed)

	fateData, _ := v2.GetFateData(&v2.FateInput{
		BirthDate: born,
		Gender:    1,
		Surname:   "张",
	})

	lastName := output.Basic().LastName
	l1 := lastName[0].ScienceStroke
	var l2 int
	if lastName[1] != nil {
		l2 = lastName[1].ScienceStroke
	}

	var nameSources []analysis.NameSource
	for _, sn := range output.AllNames() {
		nameSources = append(nameSources, analysis.NameSource{
			C1: sn.Name[0],
			C2: sn.Name[1],
		})
	}

	topNames := analysis.CollectTopNames(nameSources, "张", l1, l2, fateData, 20, func(c1, c2 *ent.Character) bool {
		if c1.CommonLevel > 2 || c2.CommonLevel > 2 {
			return false
		}
		return true
	})

	if len(topNames) == 0 {
		t.Error("CollectTopNames returned no results")
	}

	t.Logf("Top name: %s (score: %.1f, grade: %s)", topNames[0].FullName, topNames[0].Score, topNames[0].Grade)

	if topNames[0].WuGe != nil {
		t.Logf("WuGe - TianGe: %d(%s), RenGe: %d(%s), DiGe: %d(%s)",
			topNames[0].WuGe.TianGe.Stroke, topNames[0].WuGe.TianGe.Lucky,
			topNames[0].WuGe.RenGe.Stroke, topNames[0].WuGe.RenGe.Lucky,
			topNames[0].WuGe.DiGe.Stroke, topNames[0].WuGe.DiGe.Lucky)
	}

	if topNames[0].Bazi != nil {
		t.Logf("Bazi - Zodiac: %s, Constellation: %s", topNames[0].Bazi.Zodiac, topNames[0].Bazi.Constellation)
	}

	if topNames[0].ZhouYi != nil {
		t.Logf("ZhouYi - BenGua: %s(%s), Score: %d", topNames[0].ZhouYi.BenGuaName, topNames[0].ZhouYi.BenGuaJiXiong, topNames[0].ZhouYi.Score)
	}
}
