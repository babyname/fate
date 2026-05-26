// Package analysis provides name analysis and report generation.
package analysis

import (
	"fmt"
	"sort"

	v2 "github.com/godcong/chronos/v2"
	"github.com/babyname/fate/ent"
	"github.com/babyname/fate/internal/rating"
	"github.com/babyname/fate/internal/wuge"
	"github.com/babyname/fate/internal/wuxing"
)

// BuildNameResult builds a single NameResult from the given character data and fate information.
// It calculates wuge (five strokes), sancai (three talents), luck analysis, zhouyi hexagram,
// and the overall name rating score.
func BuildNameResult(rank int, surname string, c1, c2 *ent.Character, l1, l2 int, fateData *v2.FateData) NameResult {
	f1 := c1.ScienceStroke
	f2 := c2.ScienceStroke

	ge := wuge.CalcWuGe(l1, l2, f1, f2)
	tianGe := ge.TianGe()
	renGe := ge.RenGe()
	diGe := ge.DiGe()
	waiGe := ge.WaiGe()
	zongGe := ge.ZongGe()

	tianDaYan := wuge.Find(tianGe)
	renDaYan := wuge.Find(renGe)
	diDaYan := wuge.Find(diGe)
	waiDaYan := wuge.Find(waiGe)
	zongDaYan := wuge.Find(zongGe)

	sanCaiObj := wuxing.NewSanCai(tianGe, renGe, diGe)
	sanCaiStr := sanCaiObj.String()
	sanCaiLuck, _ := wuxing.GetWuXing(sanCaiStr)
	sanCaiDetail := getSanCaiDetail(sanCaiStr)

	tianWuXing := yinYangWuXingAttr(tianGe)
	renWuXing := yinYangWuXingAttr(renGe)
	diWuXing := yinYangWuXingAttr(diGe)
	waiWuXing := yinYangWuXingAttr(waiGe)
	zongWuXing := yinYangWuXingAttr(zongGe)

	renWX := string([]rune(renWuXing)[1:])
	diWX := string([]rune(diWuXing)[1:])
	tianWX := string([]rune(tianWuXing)[1:])
	waiWX := string([]rune(waiWuXing)[1:])

	jiChuYun := getJiChuYun(renWX, diWX)
	chengGongYun := getChengGongYun(renWX, tianWX)
	renJiGuanXi := getRenJiGuanXi(renWX, waiWX)

	p1 := ""
	if len(c1.Pinyin) > 0 {
		p1 = c1.Pinyin[0]
	}
	p2 := ""
	if len(c2.Pinyin) > 0 {
		p2 = c2.Pinyin[0]
	}

	isXiYong1 := false
	isXiYong2 := false
	if fateData != nil && fateData.WuxingXiji != nil {
		for _, wx := range fateData.WuxingXiji.FavorableElements {
			if c1.WuXing == wx {
				isXiYong1 = true
			}
			if c2.WuXing == wx {
				isXiYong2 = true
			}
		}
	}

	var nr *rating.NameRating
	rater := rating.NewRaterWithStrokes(fateData, l1, l2)
	nr = rater.RateName(surname, c1, c2)

	strokeStr := fmt.Sprintf("%d,%d,%d", l1, f1, f2)
	if l2 > 0 {
		strokeStr = fmt.Sprintf("%d,%d,%d,%d", l1, l2, f1, f2)
	}

	traditionalChar1 := getTraditionalChar(c1)
	traditionalChar2 := getTraditionalChar(c2)

	zhouyiResult := CalcZhouYi(l1, l2, f1, f2)
	scoreDetail := calcScoreDetail(nr)

	return NameResult{
		Rank:      rank,
		FullName:  surname + c1.Char + c2.Char,
		Surname:   surname,
		FirstName: c1.Char + c2.Char,
		Strokes:   strokeStr,
		Char1: CharInfo{
			Char:              c1.Char,
			TraditionalChar:   traditionalChar1,
			Pinyin:            p1,
			WuXing:            c1.WuXing,
			SimplifiedStroke:  c1.SimplifiedStroke,
			TraditionalStroke: c1.TraditionalStroke,
			ScienceStroke:     c1.ScienceStroke,
			KangxiStroke:      c1.KangxiStroke,
			Radical:           c1.Radical,
			Meaning:           c1.Meaning,
			IsXiYong:          isXiYong1,
		},
		Char2: CharInfo{
			Char:              c2.Char,
			TraditionalChar:   traditionalChar2,
			Pinyin:            p2,
			WuXing:            c2.WuXing,
			SimplifiedStroke:  c2.SimplifiedStroke,
			TraditionalStroke: c2.TraditionalStroke,
			ScienceStroke:     c2.ScienceStroke,
			KangxiStroke:      c2.KangxiStroke,
			Radical:           c2.Radical,
			Meaning:           c2.Meaning,
			IsXiYong:          isXiYong2,
		},
		WuGe: &WuGeResult{
			TianGe: GeItem{
				Name:          "天格",
				Stroke:        tianGe,
				Lucky:         tianDaYan.Lucky,
				DaYan:         tianDaYan.SkyNine,
				SkyNine:       tianDaYan.Comment,
				YinYangWuXing: tianWuXing,
				Analysis:      tianDaYan.Comment,
			},
			RenGe: GeItem{
				Name:          "人格",
				Stroke:        renGe,
				Lucky:         renDaYan.Lucky,
				DaYan:         renDaYan.SkyNine,
				SkyNine:       renDaYan.Comment,
				YinYangWuXing: renWuXing,
				Analysis:      renDaYan.Comment,
			},
			DiGe: GeItem{
				Name:          "地格",
				Stroke:        diGe,
				Lucky:         diDaYan.Lucky,
				DaYan:         diDaYan.SkyNine,
				SkyNine:       diDaYan.Comment,
				YinYangWuXing: diWuXing,
				Analysis:      diDaYan.Comment,
			},
			WaiGe: GeItem{
				Name:          "外格",
				Stroke:        waiGe,
				Lucky:         waiDaYan.Lucky,
				DaYan:         waiDaYan.SkyNine,
				SkyNine:       waiDaYan.Comment,
				YinYangWuXing: waiWuXing,
				Analysis:      waiDaYan.Comment,
			},
			ZongGe: GeItem{
				Name:          "总格",
				Stroke:        zongGe,
				Lucky:         zongDaYan.Lucky,
				DaYan:         zongDaYan.SkyNine,
				SkyNine:       zongDaYan.Comment,
				YinYangWuXing: zongWuXing,
				Analysis:      zongDaYan.Comment,
			},
		},
		SanCai:       sanCaiStr,
		SanCaiLuck:   sanCaiLuck,
		SanCaiDetail: sanCaiDetail,
		JiChuYun:     jiChuYun,
		ChengGongYun: chengGongYun,
		RenJiGuanXi:  renJiGuanXi,
		ZhouYi:       zhouyiResult,
		Bazi:         buildBaziSection(fateData),
		WuXing:       buildWuXingSection(fateData),
		Score:        nr.TotalScore,
		Grade:        nr.Grade,
		ScoreDetail:  scoreDetail,
		Interpret:    nr.Interpret,
	}
}

// CollectTopNames collects name results from the given sources, optionally filters them,
// sorts by score in descending order, and returns the top N results with ranks assigned.
func CollectTopNames(names []NameSource, surname string, l1, l2 int, fateData *v2.FateData, topN int, filterFunc func(c1, c2 *ent.Character) bool) []NameResult {
	var all []NameResult
	for _, nm := range names {
		c1 := nm.C1
		c2 := nm.C2

		if filterFunc != nil && !filterFunc(c1, c2) {
			continue
		}

		nr := BuildNameResult(0, surname, c1, c2, l1, l2, fateData)
		all = append(all, nr)
	}

	sort.Slice(all, func(i, j int) bool {
		return all[i].Score > all[j].Score
	})

	if topN > len(all) {
		topN = len(all)
	}

	result := all[:topN]
	for i := range result {
		result[i].Rank = i + 1
	}
	return result
}

func buildBaziSection(fateData *v2.FateData) *BaziBasic {
	if fateData == nil || fateData.Bazi == nil {
		return nil
	}
	return &BaziBasic{
		FourPillars:   fateData.Bazi.FourPillars,
		FiveElements:  fateData.Bazi.FiveElements,
		NaYin:         fateData.Bazi.NaYin,
		Zodiac:        fateData.Bazi.Zodiac,
		Constellation: fateData.Bazi.Constellation,
	}
}

func buildWuXingSection(fateData *v2.FateData) *WuXingSection {
	if fateData == nil || fateData.WuxingXiji == nil {
		return nil
	}
	ws := &WuXingSection{
		DayGan:            fateData.WuxingXiji.DayGan,
		DayWuxing:         fateData.WuxingXiji.DayWuxing,
		Strength:          fateData.WuxingXiji.Strength,
		FavorableElements: fateData.WuxingXiji.FavorableElements,
		UsefulElement:     fateData.WuxingXiji.UsefulElement,
		UnfavorableElements: fateData.WuxingXiji.UnfavorableElements,
		HostileElements:   fateData.WuxingXiji.HostileElements,
		IdleElements:      fateData.WuxingXiji.IdleElements,
		Method:            fateData.WuxingXiji.MethodName,
		MethodName:        fateData.WuxingXiji.MethodName,
		Analysis:          fateData.WuxingXiji.Analysis,
	}
	if fateData.WuxingXiji.GeJu != nil {
		ws.GeJuName = fateData.WuxingXiji.GeJu.Name
	}
	return ws
}
