package chronosfate

import (
	"fmt"

	"github.com/godcong/chronos/v2"
)

func determineGeJu(bazi BaziInfo, strength WuxingStrength) *GeJuInfo {
	siZhu := bazi.SiZhu
	riZhuGan := string([]rune(siZhu[2])[:1])
	yueZhi := string([]rune(siZhu[1])[1:])

	var mainQi string
	if stems, ok := chronos.DiZhiHiddenStems[yueZhi]; ok && len(stems) > 0 {
		mainQi = stems[0]
	}
	if mainQi == "" {
		return nil
	}

	shiShen := getShiShen(riZhuGan, mainQi)
	geJuType := shiShenToGeJu(shiShen)

	if geJuType == GeJuUnknown {
		return nil
	}

	yangRenZhi := getYangRenZhi(riZhuGan)
	isYangRen := yueZhi == yangRenZhi

	xiYongJi := geJuXiYongJi(geJuType, riZhuGan, strength, isYangRen)

	analysis := generateGeJuAnalysis(geJuType, riZhuGan, xiYongJi, isYangRen)

	return &GeJuInfo{
		Type:     geJuType,
		Name:     geJuTypeName(geJuType),
		YongShen: xiYongJi.Yong,
		XiShen:   xiYongJi.Xi,
		JiShen:   xiYongJi.Ji,
		ChouShen: xiYongJi.Chou,
		Analysis: analysis,
	}
}

func getShiShen(dayGan, target string) string {
	if m, ok := tianGanShiShenMap[dayGan]; ok {
		if ss, ok := m[target]; ok {
			return ss
		}
	}

	targetWx := wuxingOfTianGan(target)
	dayWx := wuxingOfTianGan(dayGan)
	if dayWx == "" || targetWx == "" {
		return ""
	}

	if rel, ok := wuxingRelations[dayWx]; ok {
		if result, ok := rel[targetWx]; ok {
			return result
		}
	}
	return ""
}

func shiShenToGeJu(shiShen string) GeJuType {
	switch shiShen {
	case "正官":
		return GeJuZhengGuan
	case "七杀", "偏官":
		return GeJuQiSha
	case "正财":
		return GeJuZhengCai
	case "偏财":
		return GeJuPianCai
	case "正印":
		return GeJuZhengYin
	case "偏印":
		return GeJuPianYin
	case "食神":
		return GeJuShiShen
	case "伤官":
		return GeJuShangGuan
	default:
		return GeJuUnknown
	}
}

func getYangRenZhi(gan string) string {
	m := map[string]string{
		"甲": "卯", "乙": "寅", "丙": "午", "丁": "巳",
		"戊": "午", "己": "巳", "庚": "酉", "辛": "申",
		"壬": "子", "癸": "亥",
	}
	return m[gan]
}

func geJuXiYongJi(geJu GeJuType, riZhuGan string, strength WuxingStrength, isYangRen bool) XiYongJiChou {
	switch geJu {
	case GeJuZhengGuan:
		return geJuZhengGuanYong(riZhuGan, strength, isYangRen)
	case GeJuQiSha:
		return geJuQiShaYong(riZhuGan, strength, isYangRen)
	case GeJuZhengCai:
		return geJuZhengCaiYong(riZhuGan, strength, isYangRen)
	case GeJuPianCai:
		return geJuPianCaiYong(riZhuGan, strength, isYangRen)
	case GeJuZhengYin:
		return geJuZhengYinYong(riZhuGan, strength, isYangRen)
	case GeJuPianYin:
		return geJuPianYinYong(riZhuGan, strength, isYangRen)
	case GeJuShiShen:
		return geJuShiShenYong(riZhuGan, strength, isYangRen)
	case GeJuShangGuan:
		return geJuShangGuanYong(riZhuGan, strength, isYangRen)
	default:
		return XiYongJiChou{}
	}
}

func geJuZhengGuanYong(gan string, strength WuxingStrength, isYangRen bool) XiYongJiChou {
	riWx := wuxingOfTianGan(gan)
	return XiYongJiChou{
		Yong: riWx,
		Xi:   findShengWoWuxing(riWx),
		Ji:   findKeWoWuxing(riWx),
		Chou: findWoKeWuxing(riWx),
	}
}

func geJuQiShaYong(gan string, strength WuxingStrength, isYangRen bool) XiYongJiChou {
	riWx := wuxingOfTianGan(gan)
	return XiYongJiChou{
		Yong: findKeWoWuxing(riWx),
		Xi:   riWx,
		Ji:   findShengWoWuxing(riWx),
		Chou: findWoShengWuxing(riWx),
	}
}

func geJuZhengCaiYong(gan string, strength WuxingStrength, isYangRen bool) XiYongJiChou {
	riWx := wuxingOfTianGan(gan)
	return XiYongJiChou{
		Yong: findWoShengWuxing(riWx),
		Xi:   riWx,
		Ji:   findKeWoWuxing(riWx),
		Chou: findShengWoWuxing(riWx),
	}
}

func geJuPianCaiYong(gan string, strength WuxingStrength, isYangRen bool) XiYongJiChou {
	return geJuZhengCaiYong(gan, strength, isYangRen)
}

func geJuZhengYinYong(gan string, strength WuxingStrength, isYangRen bool) XiYongJiChou {
	riWx := wuxingOfTianGan(gan)
	return XiYongJiChou{
		Yong: findShengWoWuxing(riWx),
		Xi:   riWx,
		Ji:   findWoKeWuxing(riWx),
		Chou: findWoShengWuxing(riWx),
	}
}

func geJuPianYinYong(gan string, strength WuxingStrength, isYangRen bool) XiYongJiChou {
	return geJuZhengYinYong(gan, strength, isYangRen)
}

func geJuShiShenYong(gan string, strength WuxingStrength, isYangRen bool) XiYongJiChou {
	riWx := wuxingOfTianGan(gan)
	return XiYongJiChou{
		Yong: findWoShengWuxing(riWx),
		Xi:   findShengWoWuxing(riWx),
		Ji:   findKeWoWuxing(riWx),
		Chou: riWx,
	}
}

func geJuShangGuanYong(gan string, strength WuxingStrength, isYangRen bool) XiYongJiChou {
	riWx := wuxingOfTianGan(gan)
	return XiYongJiChou{
		Yong: findWoKeWuxing(riWx),
		Xi:   findShengWoWuxing(riWx),
		Ji:   riWx,
		Chou: findWoShengWuxing(riWx),
	}
}

func geJuTypeName(t GeJuType) string {
	names := map[GeJuType]string{
		GeJuZhengGuan: "正官格",
		GeJuQiSha:     "七杀格",
		GeJuZhengCai:  "正财格",
		GeJuPianCai:   "偏财格",
		GeJuZhengYin:  "正印格",
		GeJuPianYin:   "偏印格",
		GeJuShiShen:   "食神格",
		GeJuShangGuan: "伤官格",
		GeJuUnknown:   "未知格局",
	}
	return names[t]
}

func generateGeJuAnalysis(geJu GeJuType, riZhuGan string, xiYongJi XiYongJiChou, isYangRen bool) string {
	name := geJuTypeName(geJu)
	riWx := wuxingOfTianGan(riZhuGan)

	analysis := fmt.Sprintf("日主%s，%s。", riWx, name)
	analysis += fmt.Sprintf("用神：%s，喜神：%s，忌神：%s，仇神：%s。", xiYongJi.Yong, xiYongJi.Xi, xiYongJi.Ji, xiYongJi.Chou)

	if isYangRen {
		analysis += "羊刃当权，需注意制化。"
	}

	return analysis
}

func generateBalanceAnalysis(strength WuxingStrength, xiYongJi XiYongJiChou) string {
	return fmt.Sprintf("五行总分：%.1f。用神：%s，喜神：%s，忌神：%s，仇神：%s。",
		strength.Total, xiYongJi.Yong, xiYongJi.Xi, xiYongJi.Ji, xiYongJi.Chou)
}
