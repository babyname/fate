package naming

import (
	"github.com/babyname/fate/internal/wuge"
	"github.com/babyname/fate/internal/wuxing"
)

type NameStroke struct {
	Last1  int `bson:"last_1"`
	Last2  int `bson:"last_2"`
	First1 int `bson:"first_1"`
	First2 int `bson:"first_2"`
}

type nameStroke struct {
	*NameStroke
	*wuxing.SanCai
	*wuge.WuGe
}

type SanCaiWuGe interface {
}

func (s *NameStroke) SanCaiWuGe() SanCaiWuGe {
	l1, l2, f1, f2 := s.Last1, s.Last2, s.First1, s.First2
	wuGe := wuge.CalcWuGe(l1, l2, f1, f2)
	sanCai := wuxing.NewSanCai(wuGe.TianGe(), wuGe.RenGe(), wuGe.DiGe())
	return &nameStroke{
		NameStroke: s,
		SanCai:     sanCai,
		WuGe:       wuGe,
	}
}
