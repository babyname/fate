package model

type FivePhase struct {
	Base    `xorm:"extends"`
	First   string `xorm:"varchar(2)"`
	Second  string `xorm:"varchar(2)"`
	Third   string `xorm:"varchar(2)"`
	Fortune string `xorm:"varchar(8)"`
}

func init() {
	Register(&FivePhase{})
}
func NewFivePhase(f, s, t string) FivePhase {
	return FivePhase{
		First:  f,
		Second: s,
		Third:  t,
	}
}

func (fp *FivePhase) GetFortune() string {
	db.Get(fp)
	return fp.Fortune
}

func (fp *FivePhase) Create(v ...interface{}) (int64, error) {
	i, e := db.Count(fp)
	if e == nil && i == 0 {
		return db.InsertOne(fp)
	}
	return 0, e
}

func (fp *FivePhase) Get() *FivePhase {
	db.Get(fp)
	return fp
}

//func (fp *FivePhase) CalculateFortune() string {
//	var f FivePhase
//	db.Where(fp).Get(&f)
//	return f.Fortune
//}

//
//func init() {
//	SetMigrate(CharacterFive{})
//}
//
//func FindFiveByMass(v interface{}, mass string) error {
//	return ORM().Where("mass = ?", mass).Find(v).Error
//}
//
//func FindFiveWithFirstByMass(v interface{}, fir string, mass []string) error {
//	return ORM().Where("first = ? and mass in (?)", fir, mass).Find(v).Error
//}

//func (f *CharacterFive) StringFive() string {
//	return f.First + f.Second + f.Third
//}
