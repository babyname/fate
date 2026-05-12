package filter

type CharacterFilterType int

const (
	CharacterFilterTypeDefault CharacterFilterType = iota
	CharacterFilterTypeChs
	CharacterFilterTypeCht
	CharacterFilterTypeKangxi
)

type FilterOption struct {
	CharacterFilter     bool
	CharacterFilterType CharacterFilterType
	MinStroke           int
	MaxStroke           int
	RegularFilter       bool
	DaYanFilter         bool
	WuXingFilter        bool
	SexFilter           bool
}
