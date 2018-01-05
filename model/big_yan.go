package model

//BigYan 大衍之数取天九
type BigYan struct {
	Base    `xorm:"extends"`
	Index   int    //use array index
	Goil    string //good or ill luck(吉凶),哈哈
	SkyNine string //天九(天九地十)
	Comment string
}

func init() {
	Register(&BigYan{})
}

func (by BigYan) Create(v ...interface{}) (int64, error) {
	i, e := db.Count(&by)
	if e == nil && i == 0 {
		return db.InsertOne(&by)
	}
	return 0, e
}
