package wuge

import (
	"strings"
	"sync"
)

const maxStroke = 30

type WuGeResult struct {
	FirstStroke1 int
	FirstStroke2 int
	TianGe       int
	RenGe        int
	DiGe         int
	WaiGe        int
	ZongGe       int
	ZongLucky    bool
	ZongSex      bool
	ZongMax      bool
}

var (
	luckyOnce sync.Once
	luckyTable [maxStroke + 1][maxStroke + 1][]WuGeResult
)

func GetLuckyByLastName(l1, l2 int) []WuGeResult {
	if l1 < 1 || l1 > maxStroke || l2 < 0 || l2 > maxStroke {
		return nil
	}
	luckyOnce.Do(buildLuckyTable)
	return luckyTable[l1][l2]
}

func buildLuckyTable() {
	for l1 := 1; l1 <= maxStroke; l1++ {
		for l2 := 0; l2 <= maxStroke; l2++ {
			var results []WuGeResult
			for f1 := 1; f1 <= maxStroke; f1++ {
				for f2 := 1; f2 <= maxStroke; f2++ {
					r := calcWuGeResult(l1, l2, f1, f2)
					if r.ZongLucky {
						results = append(results, r)
					}
				}
			}
			luckyTable[l1][l2] = results
		}
	}
}

func calcWuGeResult(l1, l2, f1, f2 int) WuGeResult {
	ge := CalcWuGe(l1, l2, f1, f2)
	tianGe := ge.TianGe()
	renGe := ge.RenGe()
	diGe := ge.DiGe()
	waiGe := ge.WaiGe()
	zongGe := ge.ZongGe()
	zongDaYan := Find(zongGe)

	return WuGeResult{
		FirstStroke1: f1,
		FirstStroke2: f2,
		TianGe:       tianGe,
		RenGe:        renGe,
		DiGe:         diGe,
		WaiGe:        waiGe,
		ZongGe:       zongGe,
		ZongLucky:    isLucky(zongDaYan.Lucky),
		ZongSex:      bool(zongDaYan.Sex),
		ZongMax:      zongDaYan.Max,
	}
}

func isLucky(s string) bool {
	return strings.Index(s, "吉") != -1
}
