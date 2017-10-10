package model

import "log"

type Character struct {
	Base
	IsSur      bool   ``                          //姓
	NameChar   string `gorm:"size:2"`             //字
	NameType   string `gorm:"size:2"`             //字属
	NameRoot   string `gorm:"size:128"`           //字根
	Radical    string `gorm:"size:16"`            //部首
	Strokes    int    `gorm:"not null default:0"` //笔画数
	TradName   string `gorm:"size:2"`             //繁体
	FixStrokes int    `gorm:"not null default:0"` //修正（繁体）笔画数
	Pinyin     string `gorm:"size:128"`           //拼音
	Comment    string `gorm:"size:1024"`          //备注
}

type RadicalChar struct {
	Strokes  int
	NameChar string
	Pinyin   string
	Radical  string
}

func init() {
	SetMigrate(Character{})
}

func FindByNameChar(v interface{}, nc string) error {
	return ORM().Where("name_char = ?", nc).First(v).Error
}

func FindCharactersByStrokes(v interface{}, s []int) error {
	return ORM().Where("fix_strokes in (?)", s).Find(v).Error
}

func FindCharactersByStroke(v interface{}, s int) error {
	return ORM().Where("fix_strokes = ?", s).Find(v).Error
}

func FindCharactersByStrokeBest(v interface{}, s int, best []string) error {
	if best != nil {
		return ORM().Where("fix_strokes = ?", s).Where("radical in (?)", best).Find(v).Error
	}
	return FindCharactersByStroke(v, s)
}

func FindCharactersWithFiveByStrokes(v interface{}, five string, s []int) error {
	if five == "" {
		return FindCharactersByStrokes(v, s)
	}
	return ORM().Where("fix_strokes in (?) and name_type = ? ", s, five).Find(v).Error
}

func UpdateCharacter(ch string, c RadicalChar) error {
	chr := Character{}
	ORM().Where("name_char = ?", ch).First(&chr)
	if chr.IsNil() {
		log.Println(ch, "is null")
		return nil
	}

	if chr.Radical != "" {
		log.Println("updated:", chr)
		return nil
	}
	if chr.FixStrokes == 0 {
		chr.FixStrokes = chr.Strokes
	}

	if chr.Strokes != c.Strokes {
		log.Println("strokes:", chr, c)
	}
	chr.Radical = c.Radical
	chr.Pinyin = c.Pinyin
	log.Println("updating:", chr.NameChar)
	return ORM().Save(&chr).Error
}
