package model

type Zodiac struct {
	ZodiacType ZodiacType
}

type ZodiacType int

const (
	ZODIAC_SHU ZodiacType = iota
	ZODIAC_NIU
	ZODIAC_HU
	ZODIAC_TU
	ZODIAC_LONG
	ZODIAC_SHE
	ZODIAC_MA
	ZODIAC_YANG
	ZODIAC_HOU
	ZODIAC_JI
	ZODIAC_GOU
	ZODIAC_ZHU
)

func NewZodiac(zodiacType ZodiacType) *Zodiac {
	z := new(Zodiac)
	z.ZodiacType = zodiacType
	return z
}

func (z ZodiacType) ToString() string {
	switch z {
	case ZODIAC_SHU:
		return "鼠"
	case ZODIAC_NIU:
		return "牛"
	case ZODIAC_HU:
		return "虎"
	case ZODIAC_TU:
		return "兔"
	case ZODIAC_LONG:
		return "龙"
	case ZODIAC_SHE:
		return "蛇"
	case ZODIAC_MA:
		return "马"
	case ZODIAC_YANG:
		return "羊"
	case ZODIAC_HOU:
		return "猴"
	case ZODIAC_JI:
		return "鸡"
	case ZODIAC_GOU:
		return "狗"
	case ZODIAC_ZHU:
		return "猪"
	}
	return ""
}
