package model

type FivePhase struct {
	Base    `xorm:"extends"`
	First   string `xorm:"varchar(2)"` //第一字
	Second  string `xorm:"varchar(2)"` //第二字
	Third   string `xorm:"varchar(2)"` //第三字
	Fortune string `xorm:"varchar(8)"` //吉凶
	Comment string `xorm:"text"`       //说明
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

func (fp *FivePhase) Get() (bool, error) {
	return db.Get(fp)
}

func (fp *FivePhase) UpdateOnly(cols ...string) (int64, error) {
	return db.ID(fp.Id).Cols(cols...).Update(fp)
}
func (fp *FivePhase) Sync() error {
	return Sync(fp)
}
