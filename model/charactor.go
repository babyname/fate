package model

import (
	"log"
	"strings"
)

type Character struct {
	Base
	IsSur      bool   `gorm:"not null default:false"`      //姓
	NameChar   string `gorm:"size:2 not null default: "`   //字
	NameType   string `gorm:"size:2 not null default: "`   //字属
	NameRoot   string `gorm:"size:128 not null default: "` //字根
	Radical    string `gorm:"size:16 not null default: "`  //部首
	Strokes    int    `gorm:"not null default:0"`          //笔画数
	TradName   string `gorm:"size:2 not null default:"`    //繁体
	FixStrokes int    `gorm:"not null default:0"`          //修正（繁体）笔画数
	Pinyin     string `gorm:"size:128 not null default:"`  //拼音
	Comment    string `gorm:"size:1024 not null default:"` //备注
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
	ch1 := chr
	if chr.IsNil() {
		log.Println(ch, "is null")
		return nil
	}


	if chr.Radical != "" {
		//log.Println("updated:", chr)
		return nil
	}
	if chr.FixStrokes == 0 {
		ch1.FixStrokes = chr.Strokes
	}

	if chr.FixStrokes != c.Strokes {
		log.Println("fix:", chr, c)
		ch1.FixStrokes = c.Strokes
	}

	//if chr.Strokes != c.Strokes {
	//	log.Println("strokes:", chr, c)
	//}
	if c.Radical != "" {
		ch1.Radical = c.Radical
	}

	if c.Pinyin != "" {
		ch1.Pinyin = c.Pinyin
	}

	log.Println("updating:", ch1.NameChar, ch1.FixStrokes, c.Strokes)
	return ORM().Table("characters").Where(&chr).Update(&ch1).Error
}
