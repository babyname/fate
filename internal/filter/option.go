package filter

type CharacterFilterType uint

const (
	CharacterFilterTypeDefault      CharacterFilterType = 0
	CharacterFilterTypeChs          CharacterFilterType = 1 << 0
	CharacterFilterTypeCht          CharacterFilterType = 1 << 1
	CharacterFilterTypeKangxi       CharacterFilterType = 1 << 2
	CharacterFilterTypeNameScience  CharacterFilterType = 1 << 3
	CharacterFilterTypeChsCht       CharacterFilterType = CharacterFilterTypeChs | CharacterFilterTypeCht
	CharacterFilterTypeAll          CharacterFilterType = CharacterFilterTypeChs | CharacterFilterTypeCht | CharacterFilterTypeKangxi
)

func (c CharacterFilterType) HasType(t CharacterFilterType) bool {
	if c == CharacterFilterTypeDefault {
		return t == CharacterFilterTypeNameScience
	}
	return (c & t) != 0
}

type StrokeMode int

const (
	StrokeModeScience StrokeMode = iota
	StrokeModeSimplified
	StrokeModeTraditional
	StrokeModeKangxi
)

type NameStyle int

const (
	NameStyleSimplified NameStyle = iota
	NameStyleTraditional
	NameStyleAuto
)

type FilterOption struct {
	CharacterFilter     bool
	CharacterFilterType CharacterFilterType
	MinStroke           int
	MaxStroke           int
	RegularFilter       bool

	CommonLevelFilter bool
	MinCommonLevel    int
	MaxCommonLevel    int

	WuXingFilter    bool
	PreferredWuXing []string
	AvoidWuXing     []string

	DaYanFilter bool
	SexFilter   bool

	StrokeMode StrokeMode

	NameStyle NameStyle

	GenderFilter string

	AvoidCharacters   []string
	RequireCharacters []string

	PreferredRadicals []string
	AvoidRadicals     []string

	AvoidPinyin []string

	PoetryMode       int
	XiYongMethod     string
	FilterStrictness string
}

func NewDefaultFilterOption() FilterOption {
	return FilterOption{
		CharacterFilter:     true,
		CharacterFilterType: CharacterFilterTypeDefault,
		MinStroke:           1,
		MaxStroke:           30,
		RegularFilter:       true,
		StrokeMode:          StrokeModeScience,
		NameStyle:           NameStyleAuto,
		GenderFilter:        "",
	}
}

func (fo FilterOption) WithCommonLevel(minLevel, maxLevel int) FilterOption {
	fo.CommonLevelFilter = true
	fo.MinCommonLevel = minLevel
	fo.MaxCommonLevel = maxLevel
	return fo
}

func (fo FilterOption) WithPreferredWuXing(elements ...string) FilterOption {
	fo.WuXingFilter = true
	fo.PreferredWuXing = append(fo.PreferredWuXing, elements...)
	return fo
}

func (fo FilterOption) WithAvoidWuXing(elements ...string) FilterOption {
	fo.WuXingFilter = true
	fo.AvoidWuXing = append(fo.AvoidWuXing, elements...)
	return fo
}

func (fo FilterOption) WithAvoidCharacters(chars ...string) FilterOption {
	fo.AvoidCharacters = append(fo.AvoidCharacters, chars...)
	return fo
}

func (fo FilterOption) WithPreferredRadicals(radicals ...string) FilterOption {
	fo.PreferredRadicals = append(fo.PreferredRadicals, radicals...)
	return fo
}

func (fo FilterOption) WithAvoidRadicals(radicals ...string) FilterOption {
	fo.AvoidRadicals = append(fo.AvoidRadicals, radicals...)
	return fo
}

func (fo FilterOption) WithStrokeMode(mode StrokeMode) FilterOption {
	fo.StrokeMode = mode
	return fo
}

func (fo FilterOption) WithNameStyle(style NameStyle) FilterOption {
	fo.NameStyle = style
	return fo
}

func (fo FilterOption) WithGenderFilter(gender string) FilterOption {
	fo.GenderFilter = gender
	return fo
}

func (fo FilterOption) WithRequireCharacters(chars ...string) FilterOption {
	fo.RequireCharacters = append(fo.RequireCharacters, chars...)
	return fo
}

func (fo FilterOption) WithAvoidPinyin(pinyin ...string) FilterOption {
	fo.AvoidPinyin = append(fo.AvoidPinyin, pinyin...)
	return fo
}

func (fo FilterOption) WithPoetryMode(mode int) FilterOption {
	fo.PoetryMode = mode
	return fo
}

func (fo FilterOption) WithXiYongMethod(method string) FilterOption {
	fo.XiYongMethod = method
	return fo
}

func (fo FilterOption) WithFilterStrictness(strictness string) FilterOption {
	fo.FilterStrictness = strictness
	return fo
}

func (fo FilterOption) WithCharacterFilterType(t CharacterFilterType) FilterOption {
	fo.CharacterFilterType = t
	return fo
}

func (fo FilterOption) AddCharacterFilterType(t CharacterFilterType) FilterOption {
	fo.CharacterFilterType |= t
	return fo
}
