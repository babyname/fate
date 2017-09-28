package model

type Five struct {
	Mass   string `gorm:"size:2"`
	First  string `gorm:"size:2"`
	Second string `gorm:"size:2"`
	Third  string `gorm:"size:2"`
}

func init() {
	SetMigrate(Five{})
}

func FindFiveByMass(v interface{}, mass string) error {
	return ORM().Where("mass = ?", mass).Find(v).Error
}

func FindFiveWithFirstByMass(v interface{}, fir, mass string) error {
	return ORM().Where("first = ? and mass = ?", fir, mass).Find(v).Error
}
