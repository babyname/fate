package model

type ThreeFive struct {
	Base
	SurStrokes    int
	SecondStrokes int
	ThirdStrokes  int
	TianGe        int
	RenGe         int
	DiGe          int
	WaiGe         int
	ZongGe        int
	TianCai       string `gorm:"size:2"`
	RenCai        string `gorm:"size:2"`
	DiCai         string `gorm:"size:2"`
}

func init() {
	SetMigrate(ThreeFive{})
}

func MakeTianGe(five ThreeFive) int {
	return five.SurStrokes +
		1
}

func MakeRenGe(five ThreeFive) int {
	return five.SurStrokes +
		five.SecondStrokes
}

func MakeDiGe(five ThreeFive) int {
	return five.SecondStrokes +
		five.ThirdStrokes
}

func MakeWaiGe(five ThreeFive) int {
	return five.ThirdStrokes +
		1
}

func MakeZongGe(five ThreeFive) int {
	return five.SurStrokes +
		five.SecondStrokes +
		five.ThirdStrokes
}

//1、2甲乙木，3、4丙丁火，5、6戊己土，7、8庚辛金，9、10壬癸水
func MakeSanCai(i int) string {
	switch i % 10 {
	case 1, 2:
		return "木"
	case 3, 4:
		return "火"
	case 5, 6:
		return "土"
	case 7, 8:
		return "金"
	case 9, 0:
		return "水"
	}
	return ""
}

func (five ThreeFive) InitSave() {
	five.DiGe = MakeDiGe(five)
	five.RenGe = MakeRenGe(five)
	five.TianGe = MakeTianGe(five)
	five.WaiGe = MakeWaiGe(five)
	five.ZongGe = MakeZongGe(five)
	five.RenCai = MakeSanCai(five.RenGe)
	five.TianCai = MakeSanCai(five.TianGe)
	five.DiCai = MakeSanCai(five.DiGe)
	ORM().Create(&five)
}

func FindSecondThreeFive(sur int) []ThreeFive {
	var tf []ThreeFive
	ORM().Where("sur_strokes = ?", sur).Find(&tf)
	return tf
}

func FindSecondStrokes(sur int) []int {
	var s []int
	tf := FindSecondThreeFive(sur)
	for _, v := range tf {
		s = append(s, v.SecondStrokes)
	}
	return s
}

func FindSecondStrokesByCharacter(character Character) []int {
	if character.FixStrokes != 0 {
		return FindSecondStrokes(character.FixStrokes)
	}
	return FindSecondStrokes(character.Strokes)
}

func FindThirdThreeFive(sur, sec int) []ThreeFive {
	var tf []ThreeFive
	ORM().Where("sur_strokes = ? and second_strokes = ?", sur, sec).Find(&tf)
	return tf
}

func FindThirdStrokes(sur, sec int) []int {
	var s []int
	tf := FindThirdThreeFive(sur, sec)
	for _, v := range tf {
		s = append(s, v.SecondStrokes)
	}
	return s
}

func FindStrokesWithFive(five Five) []ThreeFive {
	var tf []ThreeFive
	if ORM().Where("tian_cai = ? and ren_cai = ? and di_cai = ?", five.First, five.Second, five.Third).Find(&tf).Error != nil {
		return nil
	}
	return tf
}
