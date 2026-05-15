package wuge

// WuGe 表示五格（天格、人格、地格、外格、总格）的计算结果。
type WuGe struct {
	tianGe int
	renGe  int
	diGe   int
	waiGe  int
	zongGe int
}

// ZongGe 返回总格数理。
func (ge *WuGe) ZongGe() int {
	return ge.zongGe
}

// WaiGe 返回外格数理。
func (ge *WuGe) WaiGe() int {
	return ge.waiGe
}

// DiGe 返回地格数理。
func (ge *WuGe) DiGe() int {
	return ge.diGe
}

// RenGe 返回人格数理。
func (ge *WuGe) RenGe() int {
	return ge.renGe
}

// TianGe 返回天格数理。
func (ge *WuGe) TianGe() int {
	return ge.tianGe
}

// CalcWuGe 根据姓氏笔画（l1、l2）和名字笔画（f1、f2）计算五格数理。
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
	return l1 + l2 + f1 + f2
}
