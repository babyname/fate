package model

type ZodiacType string

const (
	ZodiacTypeMouse   ZodiacType = "鼠"
	ZodiacTypeCow     ZodiacType = "牛"
	ZodiacTypeTiger   ZodiacType = "虎"
	ZodiacTypeRabbit  ZodiacType = "兔"
	ZodiacTypeDragon  ZodiacType = "龙"
	ZodiacTypeSnake   ZodiacType = "蛇"
	ZodiacTypeHourse  ZodiacType = "马"
	ZodiacTypeSheep   ZodiacType = "羊"
	ZodiacTypeMonkey  ZodiacType = "猴"
	ZodiacTypeChicken ZodiacType = "鸡"
	ZodiacTypeDog     ZodiacType = "狗"
	ZodiacTypePig     ZodiacType = "猪"
)

var ZodiacTypes = []ZodiacType{
	ZodiacTypeMouse,
	ZodiacTypeCow,
	ZodiacTypeTiger,
	ZodiacTypeRabbit,
	ZodiacTypeDragon,
	ZodiacTypeSnake,
	ZodiacTypeHourse,
	ZodiacTypeSheep,
	ZodiacTypeMonkey,
	ZodiacTypeChicken,
	ZodiacTypeDog,
	ZodiacTypePig,
}

type Zodiac struct {
	Base       `xorm:"extends"`
	ZodiacType string
	Best       string ``
	Worst      string `xorm:"notnull default()"`
}

//
//func init() {
//	SetMigrate(Zodiac{})
//}
//
//func NewZodiac(zodiacType ZodiacType) *Zodiac {
//	z := new(Zodiac)
//	z.ZodiacType = zodiacType
//	return z
//}
//
func (z *Zodiac) Create() {
	DB().InsertOne(z)
}

//
//func FindZodiac(zodiacType ZodiacType) *Zodiac {
//	z := new(Zodiac)
//	ORM().Where("zodiac_type = ?", zodiacType).First(z)
//	return z
//}
