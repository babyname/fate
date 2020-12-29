package fate

// NameStroke ...
type NameStroke struct {
	//ID     bson.ObjectId `bson:"_id,omitempty"`
	Last1  int `bson:"last_1"`
	Last2  int `bson:"last_2"`
	First1 int `bson:"first_1"`
	First2 int `bson:"first_2"`
}

type nameStroke struct {
	*NameStroke
	*SanCai
	*WuGe
}

// SanCaiWuGe ...
type SanCaiWuGe interface {
}

//SanCaiWuGe 三才五格
func (s *NameStroke) SanCaiWuGe() SanCaiWuGe {
	l1, l2, f1, f2 := s.Last1, s.Last2, s.First1, s.First2
	wuGe := &WuGe{
		tianGe: tianGe(l1, l2, f1, f2),
		renGe:  renGe(l1, l2, f1, f2),
		diGe:   diGe(l1, l2, f1, f2),
		waiGe:  waiGe(l1, l2, f1, f2),
		zongGe: zongGe(l1, l2, f1, f2),
	}

	sanCai := &SanCai{
		tianCai:        sanCaiAttr(wuGe.TianGe()),
		tianCaiYinYang: yinYangAttr(wuGe.TianGe()),
		renCai:         sanCaiAttr(wuGe.RenGe()),
		renCaiYinYang:  yinYangAttr(wuGe.RenGe()),
		diCai:          sanCaiAttr(wuGe.DiGe()),
		diCaiYingYang:  yinYangAttr(wuGe.DiGe()),
	}

	return &nameStroke{
		NameStroke: s,
		SanCai:     sanCai,
		WuGe:       wuGe,
	}
}
