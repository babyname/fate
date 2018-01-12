package fate

import (
	"fmt"

	"github.com/godcong/fate/debug"
)

//type ThreeTalentAttribute string
//type ThreeTalentYinYang string

type ThreeTalentComposite struct {
	ThreeTalentAttribute string
	ThreeTalentYinYang   string
}

type ThreeTalent struct {
	SkyTalent    ThreeTalentComposite
	LandTalent   ThreeTalentComposite
	PersonTalent ThreeTalentComposite
}

//天人地格算三才
func NewThreeTalent(grid FiveGrid) ThreeTalent {
	return ThreeTalent{
		SkyTalent:    NewThreeTalentComposite(grid.SkyGrid),
		LandTalent:   NewThreeTalentComposite(grid.LandGrid),
		PersonTalent: NewThreeTalentComposite(grid.PersonGrid),
	}
}

func NewThreeTalentComposite(i int) ThreeTalentComposite {
	return ThreeTalentComposite{
		ThreeTalentAttribute: GenerateAttribute(i),
		ThreeTalentYinYang:   GenerateYinYang(i),
	}
}

// GenerateThreeTalent 计算字符的三才属性
// 1-2木：1为阳木，2为阴木   3-4火：3为阳火，4为阴火   5-6土：5为阳土，6为阴土   7-8金：7为阳金，8为阴金   9-10水：9为阳水，10为阴水
func GenerateAttribute(i int) string {
	var attr string
	switch i % 10 {
	case 1, 2:
		attr = "木"
	case 3, 4:
		attr = "火"
	case 5, 6:
		attr = "土"
	case 7, 8:
		attr = "金"
	case 9, 0:
		fallthrough
	default:
		attr = "水"
	}

	return attr
}

func GenerateYinYang(i int) string {
	var yy string
	if i%2 == 1 {
		yy = "阳"
	} else {
		yy = "阴"
	}
	return yy
}

func (t *ThreeTalent) PrintThreeTalent() {
	v := fmt.Sprintf("天才:%s,人才:%s,地才:%s",
		t.SkyTalent.ThreeTalentAttribute,
		t.PersonTalent.ThreeTalentAttribute,
		t.LandTalent.ThreeTalentAttribute,
	)
	debug.Println(v)
}

//
//type SanCai struct {
//	Base
//	//SurStrokes    int
//	//SecondStrokes int
//	//ThirdStrokes  int
//	TianGe        int
//	RenGe         int
//	DiGe          int
//	WaiGe         int
//	ZongGe        int
//	TianCai       string `gorm:"size:2"`
//	RenCai        string `gorm:"size:2"`
//	DiCai         string `gorm:"size:2"`
//}
//
//func init() {
//	SetMigrate(ThreeFive{})
//}
//
//func NewThreeFive(sur, sec, trd int) ThreeFive {
//	tf := ThreeFive{
//		SurStrokes:    sur,
//		SecondStrokes: sec,
//		ThirdStrokes:  trd,
//	}
//	tf.DiGe = MakeDiGe(tf)
//	tf.RenGe = MakeRenGe(tf)
//	tf.TianGe = MakeTianGe(tf)
//	tf.WaiGe = MakeWaiGe(tf)
//	tf.ZongGe = MakeZongGe(tf)
//	tf.RenCai = MakeSanCai(tf.RenGe)
//	tf.TianCai = MakeSanCai(tf.TianGe)
//	tf.DiCai = MakeSanCai(tf.DiGe)
//	return tf
//}
//
//func (tf ThreeFive) PrintString() {
//	ps := fmt.Sprintf("总格：%d，天格：%d，人格：%d，地格：%d，外格：%d", tf.ZongGe, tf.TianGe, tf.RenGe, tf.DiGe, tf.WaiGe)
//	log.Println(ps)
//}
//
//func MakeTianGe(five ThreeFive) int {
//	return five.SurStrokes +
//		1
//}
//
//func MakeRenGe(five ThreeFive) int {
//	return five.SurStrokes +
//		five.SecondStrokes
//}
//
//func MakeDiGe(five ThreeFive) int {
//	return five.SecondStrokes +
//		five.ThirdStrokes
//}
//
//func MakeWaiGe(five ThreeFive) int {
//	return five.ThirdStrokes +
//		1
//}
//
//func MakeZongGe(five ThreeFive) int {
//	return five.SurStrokes +
//		five.SecondStrokes +
//		five.ThirdStrokes
//}
//

//
//func (five ThreeFive) InitSave() {
//	five.DiGe = MakeDiGe(five)
//	five.RenGe = MakeRenGe(five)
//	five.TianGe = MakeTianGe(five)
//	five.WaiGe = MakeWaiGe(five)
//	five.ZongGe = MakeZongGe(five)
//	five.RenCai = MakeSanCai(five.RenGe)
//	five.TianCai = MakeSanCai(five.TianGe)
//	five.DiCai = MakeSanCai(five.DiGe)
//	ORM().Create(&five)
//}
//
//func FindSecondThreeFive(sur int) []ThreeFive {
//	var tf []ThreeFive
//	ORM().Where("sur_strokes = ?", sur).Find(&tf)
//	return tf
//}
//
//func FindSecondStrokes(sur int) []int {
//	var s []int
//	tf := FindSecondThreeFive(sur)
//	for _, v := range tf {
//		s = append(s, v.SecondStrokes)
//	}
//	return s
//}
//
//func FindSecondStrokesByCharacter(character Character) []int {
//	if character.FixStrokes != 0 {
//		return FindSecondStrokes(character.FixStrokes)
//	}
//	return FindSecondStrokes(character.Strokes)
//}
//
//func FindThirdThreeFive(sur, sec int) []ThreeFive {
//	var tf []ThreeFive
//	ORM().Where("sur_strokes = ? and second_strokes = ?", sur, sec).Find(&tf)
//	return tf
//}
//
//func FindThirdStrokes(sur, sec int) []int {
//	var s []int
//	tf := FindThirdThreeFive(sur, sec)
//	for _, v := range tf {
//		s = append(s, v.SecondStrokes)
//	}
//	return s
//}
//
//func FindStrokesWithFive(five Five) []ThreeFive {
//	var tf []ThreeFive
//	if ORM().Where("tian_cai = ? and ren_cai = ? and di_cai = ?", five.First, five.Second, five.Third).Find(&tf).Error != nil {
//		return nil
//	}
//	return tf
//}
