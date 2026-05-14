package fate

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/babyname/fate/config"
	"github.com/babyname/fate/ent"
	"github.com/babyname/fate/internal/analysis"
	v2 "github.com/godcong/chronos/v2"
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
	for name, ok := output.NextName(); ok; name, ok = output.NextName() {
		nameSources = append(nameSources, analysis.NameSource{
			C1: name.FirstName[0],
			C2: name.FirstName[1],
		})
	}

	topNames := analysis.CollectTopNames(nameSources, "张", l1, l2, fateData, 20, func(c1, c2 *ent.Character) bool {
		if c1.CommonLevel > 2 || c2.CommonLevel > 2 {
			return false
		}
		return true
	})

	report := analysis.NewReport("张", born.Format("2006年01月02日 15:04"), "男", fateData, total)
	report.TopNames = topNames

	outputDir := filepath.Join(".", "output")
	os.MkdirAll(outputDir, os.ModePerm)

	formatters := []analysis.Formatter{
		&analysis.TextFormatter{},
		&analysis.MarkdownFormatter{},
		&analysis.JSONFormatter{Indent: true},
	}

	for _, fmt_ := range formatters {
		filename := "张_姓名分析报告" + fmt_.Extension()
		filepath := filepath.Join(outputDir, filename)

		f, err := os.Create(filepath)
		if err != nil {
			t.Errorf("create %s error: %v", filepath, err)
			continue
		}

		err = fmt_.Format(f, report)
		f.Close()
		if err != nil {
			t.Errorf("format %s error: %v", filepath, err)
			continue
		}

		t.Logf("已生成: %s", filepath)
	}
}
