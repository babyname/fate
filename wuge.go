package fate

import (
	"github.com/go-xorm/xorm"
	"github.com/godcong/fate/data"
)

//WuGe
type WuGe struct {
	tianGe int
	renGe  int
	diGe   int
	waiGe  int
	zongGe int
}

func (wuGe *WuGe) ZongGe() int {
	return wuGe.zongGe
}

func (wuGe *WuGe) WaiGe() int {
	return wuGe.waiGe
}

func (wuGe *WuGe) DiGe() int {
	return wuGe.diGe
}

func (wuGe *WuGe) RenGe() int {
	return wuGe.renGe
}

func (wuGe *WuGe) TianGe() int {
	return wuGe.tianGe
}

//NewWuGe 计算五格
func NewWuGe(l1, l2, f1, f2 int) *WuGe {
	return &WuGe{
		tianGe: tianGe(l1, l2, f1, f2),
		renGe:  renGe(l1, l2, f1, f2),
		diGe:   diGe(l1, l2, f1, f2),
		waiGe:  waiGe(l1, l2, f1, f2),
		zongGe: zongGe(l1, l2, f1, f2),
	}
}

//tianGe input the ScienceStrokes with last name
//天格（复姓）姓的笔画相加
//天格（单姓）姓的笔画上加一
func tianGe(l1, l2, _, _ int) int {
	if l2 == 0 {
		return l1 + 1
	}
	return l1 + l2
}

//renGe input the ScienceStrokes with name
//人格（复姓）姓氏的第二字的笔画加名的第一字
//人格（复姓单名）姓的第二字加名
//人格（单姓单名）姓加名
// 人格（单姓复名）姓加名的第一字
func renGe(l1, l2, f1, _ int) int {
	//人格（复姓）姓氏的第二字的笔画加名的第一字
	//人格（复姓单名）姓的第二字加名
	if l2 != 0 {
		return l2 + f1
	} else {
		return l1 + f1
	}
}

//diGe input the ScienceStrokes with name
//地格（复姓复名，单姓复名）名字相加
//地格（复姓单名，单姓单名）名字+1
func diGe(_, _, f1, f2 int) int {
	if f2 == 0 {
		return f1 + 1
	}
	return f1 + f2
}

//waiGe input the ScienceStrokes with name
//外格（复姓单名）姓的第一字加笔画数一
//外格（复姓复名）姓的第一字和名的最后一定相加的笔画数
//外格（单姓复名）一加名的最后一个字
//外格（单姓单名）一加一
func waiGe(l1, l2, _, f2 int) (n int) {
	//单姓单名
	if l2 == 0 && f2 == 0 {
		n = 1 + 1
	}
	//单姓复名
	if l2 == 0 && f2 != 0 {
		n = 1 + f2
	}
	//复姓单名
	if l2 != 0 && f2 == 0 {
		n = l1 + 1
	}
	//复姓复名
	if l2 != 0 && f2 != 0 {
		n = l1 + f2
	}
	return n
}

//zongGe input the ScienceStrokes with name
//总格，姓加名的笔画总数  数理五行分类
func zongGe(l1, l2, f1, f2 int) int {
	return l1 + l2 + f1 + f2
}

func checkDaYan(idx int) bool {
	switch data.DaYanList[idx-1].Lucky {
	case "吉", "半吉":
		return true
	}
	return false
}

func getDaYanLucky(idx int) string {
	if idx > 1 && idx < 81 {
		return data.DaYanList[idx-1].Lucky
	}
	return ""
}

//Check 格检查
func (ge *WuGe) Check() bool {
	//ignore:tianGe
	for _, v := range []int{ge.diGe, ge.renGe, ge.waiGe, ge.zongGe} {
		if !checkDaYan(v) {
			return false
		}
	}
	return true
}

//WuGeLucky ...
type WuGeLucky struct {
	TianGe    int    `xorm:"tian_ge"`
	TianDaYan string `xorm:"tian_da_yan"`
	RenGe     int    `xorm:"ren_ge"`
	RenDaYan  string `xorm:"ren_da_yan"`
	DiGe      int    `xorm:"di_ge"`
	DiDaYan   string `xorm:"di_da_yan"`
	WaiGe     int    `xorm:"wai_ge"`
	WaiDaYan  string `xorm:"wai_da_yan"`
	ZongGe    int    `xorm:"zong_ge"`
	ZongDaYan string `xorm:"zong_da_yan"`
	ZongLucky bool   `xorm:"zong_lucky"`
}

func InsertOrUpdate(eng *xorm.Engine, lucky *WuGeLucky) (n int64, e error) {
	ltmp := new(WuGeLucky)
	b, err := eng.Where("tian_ge", lucky.TianGe).Where("ren_ge", lucky.RenGe).Where("di_ge", lucky.DiGe).Get(ltmp)
	if err != nil {
		return 0, e
	}
	log.With("lucky", lucky).Info("exist:", b)
	if !b {
		n, e = eng.InsertOne(lucky)
		return
	}
	return eng.Where("tian_ge", lucky.TianGe).Where("ren_ge", lucky.RenGe).Where("di_ge", lucky.DiGe).Update(lucky)
}

const WuGeMax = 31

func initWuGe() <-chan *WuGeLucky {
	var wuge *WuGe
	lucky := make(chan *WuGeLucky)
	l1, l2, f1, f2 := 1, 1, 1, 1
	go func() {
		for ; l1 < WuGeMax; l1++ {
			for ; l2 < WuGeMax; l2++ {
				for ; f1 < WuGeMax; f1++ {
					for ; f2 < WuGeMax; f2++ {
						wuge = NewWuGe(l1, l2, f1, f2)
						lucky <- &WuGeLucky{
							TianGe:    wuge.tianGe,
							TianDaYan: getDaYanLucky(wuge.tianGe),
							RenGe:     wuge.renGe,
							RenDaYan:  getDaYanLucky(wuge.renGe),
							DiGe:      wuge.diGe,
							DiDaYan:   getDaYanLucky(wuge.diGe),
							WaiGe:     wuge.waiGe,
							WaiDaYan:  getDaYanLucky(wuge.waiGe),
							ZongGe:    wuge.zongGe,
							ZongDaYan: getDaYanLucky(wuge.zongGe),
							ZongLucky: wuge.Check(),
						}
					}
					f2 = 1
				}
				f1 = 1
			}
			l2 = 1
		}
		lucky <- nil
	}()

	return lucky
}
