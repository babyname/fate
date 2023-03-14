package fate

// CharacterFilterType use the character query type
// ENUM(default,chs,cht,kangxi)
type CharacterFilterType int

type FilterOption struct {
	CharacterFilter     bool
	CharacterFilterType CharacterFilterType //default,chs,cht,kx
	MinCharacter        int
	MaxCharacter        int
	RegularFilter       bool
	DaYanFilter         bool
	WuXingFilter        bool
	SexFilter           bool
}
