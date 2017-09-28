package model

type Character struct {
	Base
	IsSur      bool
	NameChar   string `gorm:"size:2"`
	NameType   string `gorm:"size:2"`   //字属
	NameRoot   string `gorm:"size:128"` //字根
	Radical    string `gorm:"size:16"`  //部首
	Strokes    int    //笔画数
	FixStrokes int    //修正（繁体）笔画数
	Comment    string `gorm:"size:1024"` //备注
}

func init() {
	SetMigrate(Character{})
}

func FindByNameChar(v interface{}, nc string) error {
	return ORM().Where("name_char = ?", nc).First(v).Error
}

func FindCharactersByStrokes(v interface{}, s []int) error {
	return ORM().Where("strokes in (?)", s).Find(v).Error
}

func FindCharactersWithFiveByStrokes(v interface{}, five string, s []int) error {
	return ORM().Where("strokes in (?) and name_type = ? ", s, five).Find(v).Error
}
