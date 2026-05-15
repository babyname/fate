package analysis

import (
	"github.com/babyname/fate/internal/rating"
)

func calcScoreDetail(nr *rating.NameRating) ScoreDetail {
	return ScoreDetail{
		WenHuaYinXiang: nr.WenHuaScore,
		WuXingBaZi:     nr.WuXingScore,
		ShengXiao:      nr.ShengXiaoScore,
		WuGeShuLi:      nr.WuGeScore,
		YinYun:         nr.YinYunScore,
	}
}
