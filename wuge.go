package fate

import (
	"github.com/goextension/log"
	"github.com/google/uuid"
	"github.com/xormsharp/xorm"
)

// WuGe ...
type WuGe struct {
	tianGe int
	renGe  int
	diGe   int
	waiGe  int
	zongGe int
}

// ZongGe ...
func (ge *WuGe) ZongGe() int {
	return ge.zongGe
}

// WaiGe ...
func (ge *WuGe) WaiGe() int {
	return ge.waiGe
}

// DiGe ...
func (ge *WuGe) DiGe() int {
	return ge.diGe
}

// RenGe ...
func (ge *WuGe) RenGe() int {
	return ge.renGe
}

// TianGe ...
func (ge *WuGe) TianGe() int {
	return ge.tianGe
}

//CalcWuGe 计算五格
func CalcWuGe(l1, l2, f1, f2 int) *WuGe {
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
	}
	return l1 + f1
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
	//归1
	zg := (l1 + l2 + f1 + f2) - 1
	if zg < 0 {
		zg = zg + 81
	}
	return zg%81 + 1
}

func checkDaYan(idx int) bool {
	switch daYanList[idx-1].Lucky {
	case "吉", "半吉":
		return true
	}
	return false
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
	ID           string `xorm:"id pk"`
	LastStroke1  int    `xorm:"last_stroke_1"`
	LastStroke2  int    `xorm:"last_stroke_2"`
	FirstStroke1 int    `xorm:"first_stroke_1"`
	FirstStroke2 int    `xorm:"first_stroke_2"`
	TianGe       int    `xorm:"tian_ge"`
	TianDaYan    string `xorm:"tian_da_yan"`
	RenGe        int    `xorm:"ren_ge"`
	RenDaYan     string `xorm:"ren_da_yan"`
	DiGe         int    `xorm:"di_ge"`
	DiDaYan      string `xorm:"di_da_yan"`
	WaiGe        int    `xorm:"wai_ge"`
	WaiDaYan     string `xorm:"wai_da_yan"`
	ZongGe       int    `xorm:"zong_ge"`
	ZongDaYan    string `xorm:"zong_da_yan"`
	ZongLucky    bool   `xorm:"zong_lucky"`
	ZongSex      bool   `xorm:"zong_sex"`
	ZongMax      bool   `xorm:"zong_max"`
}

// BeforeInsert ...
func (w *WuGeLucky) BeforeInsert() {
	w.ID = uuid.Must(uuid.NewUUID()).String()
}

func countWuGeLucky(engine *xorm.Engine) (n int64, e error) {
	return engine.Table(&WuGeLucky{}).Count()
}

func insertOrUpdateWuGeLucky(engine *xorm.Engine, lucky *WuGeLucky) (n int64, e error) {
	session := engine.Where("last_stroke_1 = ?", lucky.LastStroke1).
		Where("last_stroke_2 = ?", lucky.LastStroke2).
		Where("first_stroke_1 = ?", lucky.FirstStroke1).
		Where("first_stroke_2 = ?", lucky.FirstStroke2)

	n, e = session.Clone().Count(&WuGeLucky{})
	if e != nil {
		return n, e
	}
	log.Infow("lucky", lucky)
	if n == 0 {
		n, e = engine.InsertOne(lucky)
		return
	}
	return session.Clone().Update(lucky)
}

// WuGeMax ...
const WuGeMax = 32

func initWuGe(lucky chan<- *WuGeLucky) {
	defer func() {
		close(lucky)
	}()
	var wuge *WuGe
	for l1 := 1; l1 <= WuGeMax; l1++ {
		for l2 := 0; l2 <= WuGeMax; l2++ {
			for f1 := 1; f1 <= WuGeMax; f1++ {
				for f2 := 1; f2 <= WuGeMax; f2++ {
					wuge = CalcWuGe(l1, l2, f1, f2)
					lucky <- &WuGeLucky{
						ID:           "",
						LastStroke1:  l1,
						LastStroke2:  l2,
						FirstStroke1: f1,
						FirstStroke2: f2,
						TianGe:       wuge.tianGe,
						TianDaYan:    GetDaYan(wuge.tianGe).Lucky,
						RenGe:        wuge.renGe,
						RenDaYan:     GetDaYan(wuge.renGe).Lucky,
						DiGe:         wuge.diGe,
						DiDaYan:      GetDaYan(wuge.diGe).Lucky,
						WaiGe:        wuge.waiGe,
						WaiDaYan:     GetDaYan(wuge.waiGe).Lucky,
						ZongGe:       wuge.zongGe,
						ZongDaYan:    GetDaYan(wuge.zongGe).Lucky,
						ZongLucky:    wuge.Check(),
						ZongSex:      isSex(wuge.zongGe, wuge.waiGe, wuge.renGe, wuge.diGe),
						ZongMax:      GetDaYan(wuge.zongGe).IsMax(),
					}
				}
			}
		}
	}
}

func getStroke(character *Character) int {
	if character.ScienceStroke != 0 {
		return character.ScienceStroke
	} else if character.KangXiStroke != 0 {
		return character.KangXiStroke
	} else if character.Stroke != 0 {
		return character.Stroke
	} else if character.SimpleTotalStroke != 0 {
		return character.SimpleTotalStroke
	} else if character.TraditionalTotalStroke != 0 {
		return character.TraditionalTotalStroke
	}
	return 0
}

func isSex(dys ...int) bool {
	for _, dy := range dys {
		if GetDaYan(dy).Sex {
			return true
		}
	}
	return false
}

func filterWuGe(eng *xorm.Engine, last []*Character, wg chan<- *WuGeLucky) error {
	defer func() {
		close(wg)
	}()
	l1 := getStroke(last[0])
	l2 := 0
	if len(last) == 2 {
		l2 = getStroke(last[1])
	}
	s := eng.Where("last_stroke_1 =?", l1).
		And("last_stroke_2 =?", l2).
		And("zong_lucky = ?", 1)
	rows, e := s.Rows(&WuGeLucky{})
	if e != nil {
		return e
	}
	for rows.Next() {
		var tmp WuGeLucky
		e := rows.Scan(&tmp)
		if e != nil {
			return e
		}
		wg <- &tmp
	}

	return nil
}
