package fate

import (
	"fmt"
	"testing"
	"time"

	"github.com/babyname/fate/v4/config"
)

func TestNameGenerationPerformance(t *testing.T) {
	// Performance test on real 30k-char database.
	// Tests 4 representative last names across 2 filter modes = 8 runs.
	// Each run generates up to 10000 names; total ~80s on warm cache.
	cfg := config.DefaultConfig()
	cfg.Database = config.DBConfig{
		Driver: "sqlite3",
		Name:   testDBPath,
	}

	f, err := New(cfg)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// Representative surnames: single-char, different strokes.
	lastNames := []struct {
		name [2]string
		desc string
	}{
		{[2]string{"张", ""}, "张(11画,火)"},
		{[2]string{"李", ""}, "李(7画,木)"},
		{[2]string{"王", ""}, "王(4画,土)"},
		{[2]string{"陈", ""}, "陈(16画,火)"},
	}

	born, _ := time.Parse("2006/01/02 15:04", "2024/06/15 10:30")

	// Full filter (the heaviest path — all gates active).
	fullFilter := NewFilter(FilterOption{
		CharacterFilter:     true,
		CharacterFilterType: CharacterFilterTypeDefault,
		MinStroke:           3,
		MaxStroke:           18,
		RegularFilter:       true,
		DaYanFilter:         true,
		WuXingFilter:        true,
	})

	var totalRuns, totalNames, zeroRuns int
	var totalTime time.Duration

	for _, ln := range lastNames {
		s := f.NewSessionWithFilter(fullFilter)

		input := &Input{
			Last: ln.name,
			Born: born,
			Sex:  Sex(1),
		}

		start := time.Now()
		if err := s.Start(input); err != nil {
			t.Errorf("  %s: Start() error = %v", ln.desc, err)
			continue
		}
		s.Wait()
		elapsed := time.Since(start)

		output := input.Output()
		n := output.Total()

		t.Logf("  %s: %d names in %v (%.0f/s)", ln.desc, n, elapsed, float64(n)/elapsed.Seconds())

		totalRuns++
		totalNames += n
		totalTime += elapsed
		if n == 0 {
			zeroRuns++
		}
	}

	avgTime := totalTime / time.Duration(totalRuns)
	t.Logf("=== Summary ===")
	t.Logf("  Runs: %d | Avg names: %d | Avg time: %v | Throughput: %.0f names/s",
		totalRuns, totalNames/totalRuns, avgTime, float64(totalNames/totalRuns)/avgTime.Seconds())

	if zeroRuns > 0 {
		t.Errorf("  FATAL: %d/%d runs produced 0 names (empty database?)", zeroRuns, totalRuns)
	}
}

func BenchmarkNameGeneration_Zhang(b *testing.B) {
	cfg := config.DefaultConfig()
	cfg.Database = config.DBConfig{
		Driver: "sqlite3",
		Name:   testDBPath,
	}

	f, err := New(cfg)
	if err != nil {
		b.Fatalf("New() error = %v", err)
	}

	born, _ := time.Parse("2006/01/02 15:04", "2024/06/15 10:30")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
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

		err = s.Start(input)
		if err != nil {
			b.Fatalf("Start() error = %v", err)
		}
		s.Wait()
	}
}

func TestNameGenerationDetailedTiming(t *testing.T) {
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

	runs := 5
	var totalNames int64
	var totalTime time.Duration

	for i := 0; i < runs; i++ {
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
		totalNames += int64(total)
		totalTime += elapsed

		t.Logf("  Run %d: %d names in %v", i+1, total, elapsed)
	}

	avgNames := totalNames / int64(runs)
	avgTime := totalTime / time.Duration(runs)
	t.Logf("\n========== 汇总 ==========")
	t.Logf("  姓氏: 张")
	t.Logf("  运行次数: %d", runs)
	t.Logf("  平均名字数: %d", avgNames)
	t.Logf("  平均耗时: %v", avgTime)
	t.Logf("  吞吐量: %.0f 名/秒", float64(avgNames)/avgTime.Seconds())
	t.Logf("  总名字数: %d, 总耗时: %v", totalNames, totalTime)

	fmt.Printf("\n===== 起名性能测试结果 =====\n")
	fmt.Printf("姓氏: 张 | 运行: %d 次 | 平均: %d 名/次 | 平均耗时: %v | 吞吐: %.0f 名/秒\n",
		runs, avgNames, avgTime, float64(avgNames)/avgTime.Seconds())
}
