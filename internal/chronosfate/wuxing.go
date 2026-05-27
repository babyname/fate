package chronosfate

import (
	"strings"

	"github.com/godcong/chronos/v2"
)

func calculateWuxingXiji(lunar chronos.Lunar) WuxingXijiInfo {
	eightChar := lunar.GetEightChar()
	siZhu := eightChar.FourPillars()

	riZhuGan := string([]rune(siZhu[2])[:1])
	riZhuWuxing := wuxingOfTianGan(riZhuGan)

	qiangRuo := judgeRizhuQiangRuo(eightChar)
	tiaoHou := findTiaoHouShen(riZhuGan, lunar)

	var xi, ji string
	if qiangRuo == "强" {
		xi = findKeWoWuxing(riZhuWuxing)
		ji = riZhuWuxing
	} else {
		xi = riZhuWuxing
		ji = findKeWoWuxing(riZhuWuxing)
	}

	return WuxingXijiInfo{
		Xi:            xi,
		Ji:            ji,
		RiZhuQiangRuo: qiangRuo,
		TiaoHouShen:   tiaoHou,
	}
}

func calculateWuxingStrength(lunar chronos.Lunar) WuxingStrength {
	eightChar := lunar.GetEightChar()
	wuxingFen := calculateWuXingFen(eightChar)

	total := 0.0
	for _, v := range wuxingFen {
		total += v
	}

	return WuxingStrength{
		WuxingFen: wuxingFen,
		Total:     total,
	}
}

func judgeRizhuQiangRuo(eightChar chronos.EightChar) string {
	wuxingFen := calculateWuXingFen(eightChar)
	riZhuGan := string([]rune(eightChar.FourPillars()[2])[:1])
	riZhuWuxing := wuxingOfTianGan(riZhuGan)

	tonglei := tongleiWuxing(riZhuWuxing)
	myScore := 0.0
	for wx, score := range wuxingFen {
		if tonglei[wx] {
			myScore += score
		}
	}

	if myScore > wuxingFen["total"]/2 {
		return "强"
	}
	return "弱"
}

func findTiaoHouShen(riZhuGan string, lunar chronos.Lunar) string {
	month := lunar.GetMonth()
	tiaoHouGan := getTiaoHouTianGan(riZhuGan, month)
	if tiaoHouGan == "" {
		return ""
	}
	return tiaoHouGan + "(" + wuxingOfTianGan(tiaoHouGan) + ")"
}

func calculateWuXingFen(eightChar chronos.EightChar) map[string]float64 {
	wuxingFen := map[string]float64{
		"木": 0, "火": 0, "土": 0, "金": 0, "水": 0,
	}

	siZhu := eightChar.FourPillars()
	for _, gz := range siZhu {
		runes := []rune(gz)
		if len(runes) < 2 {
			continue
		}
		gan := string(runes[:1])
		zhi := string(runes[1:])

		ganWx := wuxingOfTianGan(gan)
		wuxingFen[ganWx] += 1.0

		if entries, ok := diZhiHiddenStemsWeighted[zhi]; ok {
			for _, entry := range entries {
				wx := wuxingOfTianGan(entry.Stem)
				wuxingFen[wx] += entry.Weight
			}
		}
	}

	total := 0.0
	for _, v := range wuxingFen {
		total += v
	}
	wuxingFen["total"] = total

	return wuxingFen
}

func getTiaoHouTianGan(dayGan string, month int) string {
	tiaoHouTable := map[string]map[int]string{
		"甲": {1: "丙", 2: "丙", 3: "丙", 4: "丙", 5: "丙", 6: "丙",
			7: "丙", 8: "丙", 9: "丙", 10: "丙", 11: "丙", 12: "丙"},
		"乙": {1: "丙", 2: "丙", 3: "丙", 4: "丙", 5: "丙", 6: "丙",
			7: "丙", 8: "丙", 9: "丙", 10: "丙", 11: "丙", 12: "丙"},
		"丙": {1: "壬", 2: "壬", 3: "壬", 4: "壬", 5: "壬", 6: "壬",
			7: "壬", 8: "壬", 9: "壬", 10: "壬", 11: "壬", 12: "壬"},
		"丁": {1: "壬", 2: "壬", 3: "壬", 4: "壬", 5: "壬", 6: "壬",
			7: "壬", 8: "壬", 9: "壬", 10: "壬", 11: "壬", 12: "壬"},
		"戊": {1: "甲", 2: "甲", 3: "甲", 4: "甲", 5: "甲", 6: "甲",
			7: "甲", 8: "甲", 9: "甲", 10: "甲", 11: "甲", 12: "甲"},
		"己": {1: "甲", 2: "甲", 3: "甲", 4: "甲", 5: "甲", 6: "甲",
			7: "甲", 8: "甲", 9: "甲", 10: "甲", 11: "甲", 12: "甲"},
		"庚": {1: "丁", 2: "丁", 3: "丁", 4: "丁", 5: "丁", 6: "丁",
			7: "丁", 8: "丁", 9: "丁", 10: "丁", 11: "丁", 12: "丁"},
		"辛": {1: "丁", 2: "丁", 3: "丁", 4: "丁", 5: "丁", 6: "丁",
			7: "丁", 8: "丁", 9: "丁", 10: "丁", 11: "丁", 12: "丁"},
		"壬": {1: "戊", 2: "戊", 3: "戊", 4: "戊", 5: "戊", 6: "戊",
			7: "戊", 8: "戊", 9: "戊", 10: "戊", 11: "戊", 12: "戊"},
		"癸": {1: "戊", 2: "戊", 3: "戊", 4: "戊", 5: "戊", 6: "戊",
			7: "戊", 8: "戊", 9: "戊", 10: "戊", 11: "戊", 12: "戊"},
	}

	if m, ok := tiaoHouTable[dayGan]; ok {
		if gan, ok := m[month]; ok {
			return gan
		}
	}
	return ""
}

func wuxingOfTianGan(gan string) string {
	if wx, ok := tianGanWuxingMap[gan]; ok {
		return wx
	}
	return ""
}

func tongleiWuxing(wx string) map[string]bool {
	result := make(map[string]bool)
	switch wx {
	case "木":
		result["木"] = true
		result["水"] = true
	case "火":
		result["火"] = true
		result["木"] = true
	case "土":
		result["土"] = true
		result["火"] = true
	case "金":
		result["金"] = true
		result["土"] = true
	case "水":
		result["水"] = true
		result["金"] = true
	}
	return result
}

func findKeWoWuxing(myWuxing string) string {
	keMap := map[string]string{
		"木": "金",
		"火": "水",
		"土": "木",
		"金": "火",
		"水": "土",
	}
	if ke, ok := keMap[myWuxing]; ok {
		return ke
	}
	return ""
}

func joinStrings(ss []string, sep string) string {
	return strings.Join(ss, sep)
}
