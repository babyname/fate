package fate

type FilterOption struct {
	CharacterFilter     bool
	CharacterFilterType int //default,simple,trad,kangxi
	MinCharacter        int
	MaxCharacter        int
	DaYanFilter         bool
	WuXingFilter        bool
	SexFilter           bool
}
