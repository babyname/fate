package model

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
//func (z *Zodiac) Create() {
//	ORM().Create(z)
//}
//
//func FindZodiac(zodiacType ZodiacType) *Zodiac {
//	z := new(Zodiac)
//	ORM().Where("zodiac_type = ?", zodiacType).First(z)
//	return z
//}

//func (z ZodiacType) ToString() string {
//	switch z {
//	case ZODIAC_SHU:
//		return "鼠"
//	case ZODIAC_NIU:
//		return "牛"
//	case ZODIAC_HU:
//		return "虎"
//	case ZODIAC_TU:
//		return "兔"
//	case ZODIAC_LONG:
//		return "龙"
//	case ZODIAC_SHE:
//		return "蛇"
//	case ZODIAC_MA:
//		return "马"
//	case ZODIAC_YANG:
//		return "羊"
//	case ZODIAC_HOU:
//		return "猴"
//	case ZODIAC_JI:
//		return "鸡"
//	case ZODIAC_GOU:
//		return "狗"
//	case ZODIAC_ZHU:
//		return "猪"
//	}
//	return ""
//}
