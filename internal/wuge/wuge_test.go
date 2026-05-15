package wuge_test

import (
	"strings"
	"testing"

	"github.com/babyname/fate/internal/wuge"
)

func TestWuGe_WaiGe(t *testing.T) {
	l1, l2, f1, f2 := 1, 1, 1, 1
	for i := 0; i < 80000; i++ {
		if f2 >= 30 {
			f1++
			f2 = 1
		}
		if f1 >= 30 {
			l2++
			f1 = 1
		}
		if l2 >= 30 {
			l1++
			l2 = 1
		}
		wg := wuge.CalcWuGe(l1, l2, f1, f2)
		expected := ((l1+l2+f1+f2)-1)%81 + 1
		if wg.ZongGe() != expected {
			t.Errorf("ZongGe mismatch: l1=%d l2=%d f1=%d f2=%d got=%d want=%d",
				l1, l2, f1, f2, wg.ZongGe(), expected)
		}
		f2++
	}
}

func TestGetLuckyByLastName(t *testing.T) {
	results := wuge.GetLuckyByLastName(7, 8)
	if len(results) == 0 {
		t.Fatal("expected lucky results for l1=7, l2=8")
	}
	for _, r := range results {
		if r.ZongLucky != true {
			t.Errorf("expected ZongLucky=true, got false for f1=%d f2=%d", r.FirstStroke1, r.FirstStroke2)
		}
	}
	t.Logf("l1=7, l2=8: %d lucky combinations out of 900", len(results))
}

func TestGetLuckyByLastName_Boundary(t *testing.T) {
	if r := wuge.GetLuckyByLastName(0, 5); r != nil {
		t.Error("expected nil for l1=0")
	}
	if r := wuge.GetLuckyByLastName(31, 5); r != nil {
		t.Error("expected nil for l1=31")
	}
	if r := wuge.GetLuckyByLastName(5, -1); r != nil {
		t.Error("expected nil for l2=-1")
	}
	if r := wuge.GetLuckyByLastName(5, 31); r != nil {
		t.Error("expected nil for l2=31")
	}
	if r := wuge.GetLuckyByLastName(1, 0); len(r) == 0 {
		t.Error("expected results for l1=1, l2=0 (single character last name)")
	}
}

func TestGetLuckyByLastName_Consistency(t *testing.T) {
	l1, l2 := 10, 6
	results := wuge.GetLuckyByLastName(l1, l2)
	for _, r := range results {
		ge := wuge.CalcWuGe(l1, l2, r.FirstStroke1, r.FirstStroke2)
		if ge.TianGe() != r.TianGe {
			t.Errorf("TianGe mismatch: got=%d want=%d", r.TianGe, ge.TianGe())
		}
		if ge.RenGe() != r.RenGe {
			t.Errorf("RenGe mismatch: got=%d want=%d", r.RenGe, ge.RenGe())
		}
		if ge.DiGe() != r.DiGe {
			t.Errorf("DiGe mismatch: got=%d want=%d", r.DiGe, ge.DiGe())
		}
		if ge.WaiGe() != r.WaiGe {
			t.Errorf("WaiGe mismatch: got=%d want=%d", r.WaiGe, ge.WaiGe())
		}
		if ge.ZongGe() != r.ZongGe {
			t.Errorf("ZongGe mismatch: got=%d want=%d", r.ZongGe, ge.ZongGe())
		}
		daYan := wuge.Find(r.ZongGe)
		if isLucky(daYan.Lucky) != r.ZongLucky {
			t.Errorf("ZongLucky mismatch for f1=%d f2=%d zongGe=%d daYan.Lucky=%s",
				r.FirstStroke1, r.FirstStroke2, r.ZongGe, daYan.Lucky)
		}
	}
}

func TestGetLuckyByLastName_TotalCount(t *testing.T) {
	totalLucky := 0
	totalAll := 0
	for l1 := 1; l1 <= 30; l1++ {
		for l2 := 0; l2 <= 30; l2++ {
			results := wuge.GetLuckyByLastName(l1, l2)
			totalLucky += len(results)
			totalAll += 30 * 30
		}
	}
	t.Logf("Total lucky: %d / %d (%.1f%%)", totalLucky, totalAll, float64(totalLucky)/float64(totalAll)*100)
}

func isLucky(s string) bool {
	return strings.Contains(s, "吉")
}

func BenchmarkGetLuckyByLastName(b *testing.B) {
	_ = wuge.GetLuckyByLastName(7, 8)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = wuge.GetLuckyByLastName(7, 8)
	}
}

func BenchmarkGetLuckyByLastName_VariousStrokes(b *testing.B) {
	strokes := [][2]int{{1, 0}, {7, 8}, {15, 10}, {30, 30}, {20, 5}, {3, 12}}
	_ = wuge.GetLuckyByLastName(1, 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := strokes[i%len(strokes)]
		_ = wuge.GetLuckyByLastName(s[0], s[1])
	}
}

func BenchmarkCalcWuGePerCombination(b *testing.B) {
	l1, l2 := 7, 8
	f1, f2 := 10, 6
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ge := wuge.CalcWuGe(l1, l2, f1, f2)
		_ = wuge.Find(ge.ZongGe())
	}
}

func BenchmarkCalcAllCombinations_OldWay(b *testing.B) {
	l1, l2 := 7, 8
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for f1 := 1; f1 <= 30; f1++ {
			for f2 := 1; f2 <= 30; f2++ {
				ge := wuge.CalcWuGe(l1, l2, f1, f2)
				daYan := wuge.Find(ge.ZongGe())
				if isLucky(daYan.Lucky) {
					_ = ge
				}
			}
		}
	}
}

func BenchmarkBuildLuckyTable(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buildLuckyTableOnce()
	}
}

var benchLuckyTable [31][31][]wuge.WuGeResult

func buildLuckyTableOnce() {
	for l1 := 1; l1 <= 30; l1++ {
		for l2 := 0; l2 <= 30; l2++ {
			var results []wuge.WuGeResult
			for f1 := 1; f1 <= 30; f1++ {
				for f2 := 1; f2 <= 30; f2++ {
					ge := wuge.CalcWuGe(l1, l2, f1, f2)
					daYan := wuge.Find(ge.ZongGe())
					if isLucky(daYan.Lucky) {
						results = append(results, wuge.WuGeResult{
							FirstStroke1: f1,
							FirstStroke2: f2,
							TianGe:       ge.TianGe(),
							RenGe:        ge.RenGe(),
							DiGe:         ge.DiGe(),
							WaiGe:        ge.WaiGe(),
							ZongGe:       ge.ZongGe(),
							ZongLucky:    true,
							ZongSex:      bool(daYan.Sex),
							ZongMax:      daYan.Max,
						})
					}
				}
			}
			benchLuckyTable[l1][l2] = results
		}
	}
}
