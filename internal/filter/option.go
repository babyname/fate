package filter

// CharacterFilterType represents the character stroke counting method.
type CharacterFilterType int

const (
	// CharacterFilterTypeDefault uses name science strokes (姓名学笔画).
	CharacterFilterTypeDefault CharacterFilterType = iota
	// CharacterFilterTypeChs uses simplified Chinese strokes (简体笔画).
	CharacterFilterTypeChs
	// CharacterFilterTypeCht uses traditional Chinese strokes (繁体笔画).
	CharacterFilterTypeCht
	// CharacterFilterTypeKangxi uses Kangxi dictionary strokes (康熙笔画).
	CharacterFilterTypeKangxi
)

// StrokeMode determines which stroke field is used for WuGe (五格) calculation.
type StrokeMode int

const (
	// StrokeModeScience uses ScienceStroke field (姓名学笔画), the default mode.
	StrokeModeScience StrokeMode = iota
	// StrokeModeSimplified uses SimplifiedStroke field (简体笔画).
	StrokeModeSimplified
	// StrokeModeTraditional uses TraditionalStroke field (繁体笔画).
	StrokeModeTraditional
	// StrokeModeKangxi uses KangxiStroke field (康熙笔画).
	StrokeModeKangxi
)

// NameStyle controls the output form of generated names.
type NameStyle int

const (
	// NameStyleSimplified outputs names in simplified Chinese (简体输出).
	NameStyleSimplified NameStyle = iota
	// NameStyleTraditional outputs names in traditional Chinese (繁体输出).
	NameStyleTraditional
)

// FilterOption contains all customizable options for name generation.
//
//nolint:revive // stuttering type name is acceptable here for API clarity
type FilterOption struct {
	// Character filtering options
	CharacterFilter     bool
	CharacterFilterType CharacterFilterType
	MinStroke           int
	MaxStroke           int
	RegularFilter       bool // Use regular characters only

	// Common character options
	CommonLevelFilter bool // Filter by common level
	MinCommonLevel    int  // Minimum common level (1-5)
	MaxCommonLevel    int  // Maximum common level (1-5)

	// Wu Xing options
	WuXingFilter    bool
	PreferredWuXing []string // Preferred Wu Xing elements (木, 火, 土, 金, 水)
	AvoidWuXing     []string // Wu Xing elements to avoid

	// Da Yan and Geomancy options
	DaYanFilter bool
	SexFilter   bool // Gender-specific lucky strokes

	// Stroke mode for WuGe calculation
	StrokeMode StrokeMode // Which stroke counting method to use

	// Name output style
	NameStyle NameStyle // Output name in simplified or traditional

	// Gender filter for name generation
	GenderFilter string // "male", "female", "neutral", or "" for all

	// Character exclusion/inclusion
	AvoidCharacters   []string // Characters to avoid
	RequireCharacters []string // Characters that must be included (at least one)

	// Radical options
	PreferredRadicals []string // Preferred radicals (部首)
	AvoidRadicals     []string // Radicals to avoid

	// Pinyin options
	AvoidPinyin []string // Pinyin to avoid (for better sound)

	// Poetry options
	PoetryMode    int // 0=off, 1=prefer, 2=only
	XiYongMethod  string
	FilterStrictness string
}

// NewDefaultFilterOption creates a FilterOption with sensible defaults.
func NewDefaultFilterOption() FilterOption {
	return FilterOption{
		CharacterFilter:     true,
		CharacterFilterType: CharacterFilterTypeDefault,
		MinStroke:           1,
		MaxStroke:           30,
		RegularFilter:       true,
		StrokeMode:          StrokeModeScience,
		NameStyle:           NameStyleSimplified,
		GenderFilter:        "",
	}
}

// WithCommonLevel sets the common character level filter.
func (fo FilterOption) WithCommonLevel(minLevel, maxLevel int) FilterOption {
	fo.CommonLevelFilter = true
	fo.MinCommonLevel = minLevel
	fo.MaxCommonLevel = maxLevel
	return fo
}

// WithPreferredWuXing adds preferred Wu Xing elements.
func (fo FilterOption) WithPreferredWuXing(elements ...string) FilterOption {
	fo.WuXingFilter = true
	fo.PreferredWuXing = append(fo.PreferredWuXing, elements...)
	return fo
}

// WithAvoidWuXing adds Wu Xing elements to avoid.
func (fo FilterOption) WithAvoidWuXing(elements ...string) FilterOption {
	fo.WuXingFilter = true
	fo.AvoidWuXing = append(fo.AvoidWuXing, elements...)
	return fo
}

// WithAvoidCharacters adds characters to avoid.
func (fo FilterOption) WithAvoidCharacters(chars ...string) FilterOption {
	fo.AvoidCharacters = append(fo.AvoidCharacters, chars...)
	return fo
}

// WithPreferredRadicals adds preferred radicals.
func (fo FilterOption) WithPreferredRadicals(radicals ...string) FilterOption {
	fo.PreferredRadicals = append(fo.PreferredRadicals, radicals...)
	return fo
}

// WithAvoidRadicals adds radicals to avoid.
func (fo FilterOption) WithAvoidRadicals(radicals ...string) FilterOption {
	fo.AvoidRadicals = append(fo.AvoidRadicals, radicals...)
	return fo
}

// WithStrokeMode sets the stroke counting mode for WuGe calculation.
func (fo FilterOption) WithStrokeMode(mode StrokeMode) FilterOption {
	fo.StrokeMode = mode
	return fo
}

// WithNameStyle sets the output name style (simplified or traditional).
func (fo FilterOption) WithNameStyle(style NameStyle) FilterOption {
	fo.NameStyle = style
	return fo
}

// WithGenderFilter sets the gender filter for name generation.
// Accepts "male", "female", "neutral", or "" for all genders.
func (fo FilterOption) WithGenderFilter(gender string) FilterOption {
	fo.GenderFilter = gender
	return fo
}

// WithRequireCharacters adds characters that must appear in the generated name.
func (fo FilterOption) WithRequireCharacters(chars ...string) FilterOption {
	fo.RequireCharacters = append(fo.RequireCharacters, chars...)
	return fo
}

// WithAvoidPinyin adds pinyin to avoid for better phonetic results.
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
