package fate

import (
	"fmt"
	"testing"
	"time"

	"github.com/babyname/fate/config"
)

func TestNameGenerationPerformance(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Database = config.DBConfig{
		Driver: "sqlite3",
		Name:   testDBPath,
	}

	f, err := New(cfg)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	lastNames := []struct {
		name [2]string
		desc string
	}{
		{[2]string{"张", ""}, "张(单姓,11画)"},
		{[2]string{"李", ""}, "李(单姓,7画)"},
		{[2]string{"王", ""}, "王(单姓,4画)"},
		{[2]string{"刘", ""}, "刘(单姓,6画)"},
		{[2]string{"陈", ""}, "陈(单姓,7画)"},
	}

	born, _ := time.Parse("2006/01/02 15:04", "2024/06/15 10:30")

	filterOpts := []struct {
		name   string
		option FilterOption
	}{
		{
			"基础过滤(笔画+常规)",
			FilterOption{
				CharacterFilter:     true,
				CharacterFilterType: CharacterFilterTypeDefault,
				MinStroke:           3,
				MaxStroke:           18,
				RegularFilter:       true,
			},
		},
		{
			"全过滤(笔画+常规+大衍+五行)",
			FilterOption{
				CharacterFilter:     true,
				CharacterFilterType: CharacterFilterTypeDefault,
				MinStroke:           3,
				MaxStroke:           18,
				RegularFilter:       true,
				DaYanFilter:         true,
				WuXingFilter:        true,
			},
		},
	}

	for _, fo := range filterOpts {
		t.Logf("========== 过滤模式: %s ==========", fo.name)
		for _, ln := range lastNames {
			s := f.NewSessionWithFilter(NewFilter(fo.option))

			input := &Input{
				Last: ln.name,
				Born: born,
				Sex:  Sex(1),
			}

			start := time.Now()
			err = s.Start(input)
			if err != nil {
				t.Errorf("  %s: Start() error = %v", ln.desc, err)
				continue
			}
			s.Wait()
			elapsed := time.Since(start)

			output := input.Output()
			total := output.Total()
			t.Logf("  %s: %d 个名字, 耗时 %v (%.0f 名/秒)",
				ln.desc, total, elapsed,
				float64(total)/elapsed.Seconds())
		}
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
