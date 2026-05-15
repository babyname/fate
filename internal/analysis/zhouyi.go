package analysis

import (
	"fmt"

	"github.com/godcong/yi"
)

// ZhouYiResult 表示周易卦象计算的结果
type ZhouYiResult struct {
	BenGuaName     string `json:"ben_gua_name"`
	BenGuaDesc     string `json:"ben_gua_desc"`
	BenGuaJiXiong  string `json:"ben_gua_ji_xiong"`
	BianGuaName    string `json:"bian_gua_name"`
	DongYaoDesc    string `json:"dong_yao_desc"`
	DongYaoJiXiong string `json:"dong_yao_ji_xiong"`
	DaXiang        string `json:"da_xiang"`
	YunShi         string `json:"yun_shi"`
	ShiYe          string `json:"shi_ye"`
	JingShang      string `json:"jing_shang"`
	QiuMing        string `json:"qiu_ming"`
	HunLian        string `json:"hun_lian"`
	JueCe          string `json:"jue_ce"`
	Score          int    `json:"score"`
}

// CalcZhouYi 根据姓名笔画数计算周易卦象结果
func CalcZhouYi(l1, l2, f1, f2 int) *ZhouYiResult {
	shang := (l1 + l2 + f1) % 8
	if shang == 0 {
		shang = 8
	}
	xia := (f1 + f2) % 8
	if xia == 0 {
		xia = 8
	}
	bianYao := (l1 + l2 + f1 + f2) % 6
	if bianYao == 0 {
		bianYao = 6
	}

	y := yi.NumberQiGua(xia, shang, bianYao)
	if y == nil {
		return nil
	}

	result := &ZhouYiResult{}

	benGua := y.Get(yi.BenGua)
	bianGua := y.Get(yi.BianGua)

	if benGua != nil {
		result.BenGuaName = benGua.GuaMing
		result.BenGuaDesc = benGua.GuaYi
		result.BenGuaJiXiong = benGua.JiXiong
	}

	if benGua != nil && benGua.XiangYue != "" {
		result.DaXiang = benGua.XiangYue
	}

	if bianGua != nil {
		result.BianGuaName = bianGua.GuaMing
	}

	bianYaoPos := y.BianYao()
	if bianGua != nil && bianYaoPos >= 0 && bianYaoPos < 6 {
		yaoNames := []string{"初爻", "二爻", "三爻", "四爻", "五爻", "上爻"}
		result.DongYaoDesc = fmt.Sprintf("动爻%s（在第%d爻）", yaoNames[bianYaoPos], bianYaoPos+1)
		yaoJiXiong := getYaoJiXiong(bianGua, bianYaoPos)
		result.DongYaoJiXiong = yaoJiXiong
	}

	if benGua != nil {
		detail := getGuaDetail(benGua.GuaMing)
		if detail != nil {
			if result.DaXiang == "" {
				result.DaXiang = detail.DaXiang
			}
			result.YunShi = detail.YunShi
			result.ShiYe = detail.ShiYe
			result.JingShang = detail.JingShang
			result.QiuMing = detail.QiuMing
			result.HunLian = detail.HunLian
			result.JueCe = detail.JueCe
		}
	}

	if result.YunShi == "" && benGua != nil {
		result.YunShi = fmt.Sprintf("本卦%s（%s），变卦%s", benGua.GuaMing, benGua.JiXiong, result.BianGuaName)
	}

	result.Score = calcZhouYiScore(result)

	return result
}

func getYaoJiXiong(gx *yi.GuaXiang, pos int) string {
	switch pos {
	case 0:
		return gx.ChuYaoJiXiong
	case 1:
		return gx.ErYaoJiXiong
	case 2:
		return gx.SanYaoJiXiong
	case 3:
		return gx.SiYaoJiXiong
	case 4:
		return gx.WuYaoJiXiong
	case 5:
		return gx.ShangYaoJiXiong
	default:
		return ""
	}
}

func calcZhouYiScore(r *ZhouYiResult) int {
	score := 70
	switch r.BenGuaJiXiong {
	case "吉":
		score += 20
	case "半吉":
		score += 10
	case "凶":
		score -= 20
	}
	switch r.DongYaoJiXiong {
	case "吉":
		score += 10
	case "半吉":
		score += 5
	case "凶":
		score -= 10
	case "平":
		score += 0
	}
	if score > 100 {
		score = 100
	}
	if score < 0 {
		score = 0
	}
	return score
}
