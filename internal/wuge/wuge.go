package wuge

type WuGe struct {
	tianGe int
	renGe  int
	diGe   int
	waiGe  int
	zongGe int
}

func (ge *WuGe) ZongGe() int {
	return ge.zongGe
}

func (ge *WuGe) WaiGe() int {
	return ge.waiGe
}

func (ge *WuGe) DiGe() int {
	return ge.diGe
}

func (ge *WuGe) RenGe() int {
	return ge.renGe
}

func (ge *WuGe) TianGe() int {
	return ge.tianGe
}

func CalcWuGe(l1, l2, f1, f2 int) *WuGe {
	return &WuGe{
		tianGe: tianGe(l1, l2, f1, f2),
		renGe:  renGe(l1, l2, f1, f2),
		diGe:   diGe(l1, l2, f1, f2),
		waiGe:  waiGe(l1, l2, f1, f2),
		zongGe: zongGe(l1, l2, f1, f2),
	}
}

func tianGe(l1, l2, _, _ int) int {
	if l2 == 0 {
		return l1 + 1
	}
	return l1 + l2
}

func renGe(l1, l2, f1, _ int) int {
	if l2 != 0 {
		return l2 + f1
	}
	return l1 + f1
}

func diGe(_, _, f1, f2 int) int {
	if f2 == 0 {
		return f1 + 1
	}
	return f1 + f2
}

func waiGe(l1, l2, _, f2 int) (n int) {
	if l2 == 0 && f2 == 0 {
		n = 1 + 1
	}
	if l2 == 0 && f2 != 0 {
		n = 1 + f2
	}
	if l2 != 0 && f2 == 0 {
		n = l1 + 1
	}
	if l2 != 0 && f2 != 0 {
		n = l1 + f2
	}
	return n
}

func zongGe(l1, l2, f1, f2 int) int {
	zg := (l1 + l2 + f1 + f2) - 1
	if zg < 0 {
		zg = zg + 81
	}
	return zg%81 + 1
}
