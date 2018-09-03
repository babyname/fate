package fate

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

//SanCaiWuGe 三才五格
func (s *NameStroke) SanCaiWuGe() *nameStroke {
	l1, l2, f1, f2 := s.Last1, s.Last2, s.First1, s.First2
	wuGe := &WuGe{
		TianGe: tianGe(l1, l2, f1, f2),
		RenGe:  renGe(l1, l2, f1, f2),
		DiGe:   diGe(l1, l2, f1, f2),
		WaiGe:  waiGe(l1, l2, f1, f2),
		ZongGe: zongGe(l1, l2, f1, f2),
	}

	sanCai := &SanCai{
		TianCai:        sanCaiAttr(wuGe.TianGe),
		TianCaiYinYang: sanCaiYinYang(wuGe.TianGe),
		RenCai:         sanCaiAttr(wuGe.RenGe),
		RenCaiYinYang:  sanCaiYinYang(wuGe.RenGe),
		DiCai:          sanCaiAttr(wuGe.DiGe),
		DiCaiYingYang:  sanCaiYinYang(wuGe.DiGe),
	}

	return &nameStroke{
		NameStroke: s,
		SanCai:     sanCai,
		WuGe:       wuGe,
	}
}
