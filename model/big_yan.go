package model

import "log"

//BigYan 大衍之数取天九
type BigYan struct {
	Base    `xorm:"extends"`
	Index   int    //use array index
	Goil    string //good or ill luck(吉凶),哈哈
	SkyNine string //天九(天九地十)
	Comment string
}

var inBigYan []BigYan

func init() {
	Register(&BigYan{})
}

func (by *BigYan) Create(v ...interface{}) (int64, error) {
	i, e := DB().Count(by)
	if e == nil && i == 0 {
		return DB().InsertOne(by)
	}
	return 0, e
}

func FindBigYanAll(by *[]BigYan) error {
	return DB().Asc("index").Find(by)
}

func FindBigYanByGoil(by *[]BigYan, v string) error {
	return DB().Find(by, BigYan{Goil: v})
}

func FindBigYanByIndex(by *[]BigYan, i int) error {
	return DB().Find(by, BigYan{Index: i})
}

func GetBigYanByIndex(i int) BigYan {
	if len(inBigYan) == 0 {
		err := FindBigYanAll(&inBigYan)
		if err != nil {
			log.Println(err)
			return BigYan{}
		}
	}

	if i == 0 || i >= len(inBigYan) {
		return BigYan{}
	}

	return inBigYan[i-1]
}
