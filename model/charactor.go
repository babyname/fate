package model

import "errors"

type Character struct {
	Base           `xorm:"extends"`
	IsSur          bool   `xorm:"notnull"`                        //姓
	SimpleChar     string `xorm:"varchar(2) index notnull"`       //字
	SimpleStrokes  int    `xorm:"notnull default(0)"`             //标准笔画数
	TradChar       string `xorm:"varchar(2) notnull default()"`   //繁体
	TradStrokes    int    `xorm:"notnull default(0)"`             //繁体笔画数
	NameType       string `xorm:"varchar(2) notnull default()"`   //五行
	NameRoot       string `xorm:"varchar(128) notnull default()"` //字根
	Radical        string `xorm:"varchar(16) notnull default()"`  //部首
	ScienceStrokes int    `xorm:"notnull default(0)"`             //姓名学笔画数
	Pinyin         string `xorm:"varchar(128) notnull"`           //拼音
	Comment        string `xorm:"text notnull default() "`        //备注
}

//type RadicalChar struct {
//	Strokes  int
//	NameChar string
//	Pinyin   string
//	Radical  string
//}

func init() {
	Register(&Character{})
}

func (c *Character) Sync() error {
	return db.Sync2(c)
}

func (c *Character) Create(v ...interface{}) (int64, error) {
	i, e := db.Count(c)
	if e == nil && i == 0 {
		return db.InsertOne(c)
	}
	return 0, e
}

func (c *Character) Get() *Character {
	db.Get(c)
	return c
}

func (c *Character) Update(v ...interface{}) (int64, error) {
	return db.Id(c.Id).Update(c)
}

func CharacterList(t string, i int, v interface{}) error {
	if t == "" && i == 0 {
		return errors.New("wrong input")
	}

	if t == "" {
		return db.Where("science_strokes = ?", i).Find(v)
	}
	if i == 0 {
		return db.Where("name_type = ?", t).Find(v)
	}
	return db.Where("name_type = ?", t).And("science_strokes = ?", i).Find(v)
}

//
//func FindByNameChar(v interface{}, nc string) error {
//	return ORM().Where("name_char = ?", nc).First(v).Error
//}
//
//func FindCharactersByStrokes(v interface{}, s []int) error {
//	return ORM().Where("fix_strokes in (?)", s).Find(v).Error
//}
//
//func FindCharactersByStroke(v interface{}, s int) error {
//	return ORM().Where("fix_strokes = ?", s).Find(v).Error
//}
//
//func FindCharactersByStrokeBest(v interface{}, s int, best []string) error {
//	if best != nil {
//		return ORM().Where("fix_strokes = ?", s).Where("radical in (?)", best).Find(v).Error
//	}
//	return FindCharactersByStroke(v, s)
//}
//
//func FindCharactersWithFiveByStrokes(v interface{}, five string, s []int) error {
//	if five == "" {
//		return FindCharactersByStrokes(v, s)
//	}
//	return ORM().Where("fix_strokes in (?) and name_type = ? ", s, five).Find(v).Error
//}

//
//func UpdateCharacter(ch string, c RadicalChar) error {
//	chr := Character{}
//	ORM().Where("name_char = ?", ch).First(&chr)
//	ch1 := chr
//	if chr.IsNil() {
//		log.Println(ch, "is null")
//		return nil
//	}
//
//	if chr.Radical != "" {
//		//log.Println("updated:", chr)
//		return nil
//	}
//	if chr.FixStrokes == 0 {
//		ch1.FixStrokes = chr.Strokes
//	}
//
//	if chr.FixStrokes != c.Strokes {
//		log.Println("fix:", chr, c)
//		ch1.FixStrokes = c.Strokes
//	}
//
//	//if chr.Strokes != c.Strokes {
//	//	log.Println("strokes:", chr, c)
//	//}
//	if c.Radical != "" {
//		ch1.Radical = c.Radical
//	}
//
//	if c.Pinyin != "" {
//		ch1.Pinyin = c.Pinyin
//	}
//
//	log.Println("updating:", ch1.NameChar, ch1.FixStrokes, c.Strokes)
//	return ORM().Table("characters").Where(&chr).Update(&ch1).Error
//}
