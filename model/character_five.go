package model

type CharacterFive struct {
	Mass   string `gorm:"size:2"`
	First  string `gorm:"size:2"`
	Second string `gorm:"size:2"`
	Third  string `gorm:"size:2"`
}

func init() {
	SetMigrate(CharacterFive{})
}

func FindFiveByMass(v interface{}, mass string) error {
	return ORM().Where("mass = ?", mass).Find(v).Error
}

func FindFiveWithFirstByMass(v interface{}, fir string, mass []string) error {
	return ORM().Where("first = ? and mass in (?)", fir, mass).Find(v).Error
}

func (f *CharacterFive) StringFive() string {
	return f.First + f.Second + f.Third
}
