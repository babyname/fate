package fate

import "github.com/godcong/yi"

// GuaYao ...
type GuaYao struct {
	Yao     string `bson:"er_yao"`          //二爻
	JiXiong string `bson:"er_yao_ji_xiong"` //二爻吉凶
}

func getYao(xiang *yi.GuaXiang, yao int) GuaYao {
	switch yao {
	case 0:
		return GuaYao{Yao: xiang.ChuYao, JiXiong: xiang.ChuYaoJiXiong}
	case 1:
		return GuaYao{Yao: xiang.ErYao, JiXiong: xiang.ErYaoJiXiong}
	case 2:
		return GuaYao{Yao: xiang.SanYao, JiXiong: xiang.SanYaoJiXiong}
	case 3:
		return GuaYao{Yao: xiang.SiYao, JiXiong: xiang.SiYaoJiXiong}
	case 4:
		return GuaYao{Yao: xiang.WuYao, JiXiong: xiang.WuYaoJiXiong}
	case 5:
		return GuaYao{Yao: xiang.ShangYao, JiXiong: xiang.ShangYaoJiXiong}
	default:
		panic("wrong yao")
	}
}

func filterYao(y *yi.Yi, fs ...string) bool {
	yao := getYao(y.Get(yi.BianGua), y.BianYao())
	for _, s := range fs {
		if yao.JiXiong == s {
			return false
		}
	}
	return true
}
