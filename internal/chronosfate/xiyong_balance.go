package chronosfate

func balanceXiYongJi(strength WuxingStrength, bazi BaziInfo) XiYongJiChou {
	weakest := findWeakestWuxing(strength.WuxingFen)
	xian := findXianWuxing(strength.WuxingFen)

	riZhuGan := string([]rune(bazi.SiZhu[2])[:1])
	riZhuWuxing := wuxingOfTianGan(riZhuGan)

	qiangRuo := "弱"
	tonglei := tongleiWuxing(riZhuWuxing)
	myScore := 0.0
	for wx, score := range strength.WuxingFen {
		if tonglei[wx] {
			myScore += score
		}
	}
	if myScore > strength.WuxingFen["total"]/2 {
		qiangRuo = "强"
	}

	var xi, yong, ji, chou string
	if qiangRuo == "强" {
		yong = findKeWoWuxing(riZhuWuxing)
		xi = findWoShengWuxing(riZhuWuxing)
		ji = riZhuWuxing
		chou = findShengWoWuxing(riZhuWuxing)
	} else {
		yong = riZhuWuxing
		xi = findShengWoWuxing(riZhuWuxing)
		ji = findKeWoWuxing(riZhuWuxing)
		chou = findWoKeWuxing(riZhuWuxing)
	}

	if weakest != "" {
		if qiangRuo == "弱" && weakest != riZhuWuxing {
			yong = weakest
		}
	}
	if xian != "" {
		if qiangRuo == "强" {
			chou = xian
		}
	}

	return XiYongJiChou{
		Xi:   xi,
		Yong: yong,
		Ji:   ji,
		Chou: chou,
	}
}

func findWeakestWuxing(wuxingFen map[string]float64) string {
	min := -1.0
	result := ""
	for _, wx := range []string{"木", "火", "土", "金", "水"} {
		score := wuxingFen[wx]
		if min < 0 || score < min {
			min = score
			result = wx
		}
	}
	return result
}

func findXianWuxing(wuxingFen map[string]float64) string {
	for _, wx := range []string{"木", "火", "土", "金", "水"} {
		if wuxingFen[wx] == 0 {
			return wx
		}
	}
	return ""
}

func findWoShengWuxing(wx string) string {
	m := map[string]string{"木": "火", "火": "土", "土": "金", "金": "水", "水": "木"}
	return m[wx]
}

func findShengWoWuxing(wx string) string {
	m := map[string]string{"木": "水", "火": "木", "土": "火", "金": "土", "水": "金"}
	return m[wx]
}

func findWoKeWuxing(wx string) string {
	m := map[string]string{"木": "土", "火": "金", "土": "水", "金": "木", "水": "火"}
	return m[wx]
}
