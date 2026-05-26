package analysis

import (
	"fmt"

	"github.com/godcong/yi"
)

func CalcZhouYi(l1, l2, f1, f2 int) *ZhouYiBasic {
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

	y := yi.DivineByNumber(shang, xia, bianYao)
	if y == nil {
		return nil
	}

	result := &ZhouYiBasic{}

	benGua := y.GetGua(yi.Ben)
	bianGua := y.GetGua(yi.Bian)

	if benGua != nil {
		result.BenGuaName = benGua.Ming
		result.BenGuaDesc = benGua.GuaYi
		result.BenGuaJiXiong = benGua.JiXiong
	}

	if benGua != nil && benGua.XiangText != "" {
		result.DaXiang = benGua.XiangText
	}

	if bianGua != nil {
		result.BianGuaName = bianGua.Ming
	}

	bianYaoPos := y.GetBianYao()
	if bianGua != nil && bianYaoPos >= 0 && bianYaoPos < 6 {
		yaoNames := []string{"初爻", "二爻", "三爻", "四爻", "五爻", "上爻"}
		result.DongYaoDesc = fmt.Sprintf("动爻%s（在第%d爻）", yaoNames[bianYaoPos], bianYaoPos+1)
		yaoJiXiong := getYaoJiXiong(bianGua, bianYaoPos)
		result.DongYaoJiXiong = yaoJiXiong
	}

	if benGua != nil {
		if result.DaXiang == "" {
			result.DaXiang = benGua.XiangText
		}
	}

	result.Score = calcZhouYiScore(result)

	return result
}

func getYaoJiXiong(gx *yi.Gua, pos int) string {
	if pos < 0 || pos >= 6 || gx.Yaos[pos] == nil {
		return ""
	}
	return gx.Yaos[pos].JiXiong
}

func calcZhouYiScore(r *ZhouYiBasic) int {
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
