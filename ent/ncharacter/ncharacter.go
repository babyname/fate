// Code generated by ent, DO NOT EDIT.

package ncharacter

const (
	// Label holds the string label denoting the ncharacter type in the database.
	Label = "ncharacter"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldPinYin holds the string denoting the pin_yin field in the database.
	FieldPinYin = "pin_yin"
	// FieldChar holds the string denoting the char field in the database.
	FieldChar = "char"
	// FieldCharStroke holds the string denoting the char_stroke field in the database.
	FieldCharStroke = "char_stroke"
	// FieldRadical holds the string denoting the radical field in the database.
	FieldRadical = "radical"
	// FieldRadicalStroke holds the string denoting the radical_stroke field in the database.
	FieldRadicalStroke = "radical_stroke"
	// FieldIsRegular holds the string denoting the is_regular field in the database.
	FieldIsRegular = "is_regular"
	// FieldIsSimplified holds the string denoting the is_simplified field in the database.
	FieldIsSimplified = "is_simplified"
	// FieldSimplifiedID holds the string denoting the simplified_id field in the database.
	FieldSimplifiedID = "simplified_id"
	// FieldIsTraditional holds the string denoting the is_traditional field in the database.
	FieldIsTraditional = "is_traditional"
	// FieldTraditionalID holds the string denoting the traditional_id field in the database.
	FieldTraditionalID = "traditional_id"
	// FieldIsKangXi holds the string denoting the is_kang_xi field in the database.
	FieldIsKangXi = "is_kang_xi"
	// FieldKangXiID holds the string denoting the kang_xi_id field in the database.
	FieldKangXiID = "kang_xi_id"
	// FieldKangXiStroke holds the string denoting the kang_xi_stroke field in the database.
	FieldKangXiStroke = "kang_xi_stroke"
	// FieldIsVariant holds the string denoting the is_variant field in the database.
	FieldIsVariant = "is_variant"
	// FieldVariantID holds the string denoting the variant_id field in the database.
	FieldVariantID = "variant_id"
	// FieldIsScience holds the string denoting the is_science field in the database.
	FieldIsScience = "is_science"
	// FieldScienceStroke holds the string denoting the science_stroke field in the database.
	FieldScienceStroke = "science_stroke"
	// FieldWuXing holds the string denoting the wu_xing field in the database.
	FieldWuXing = "wu_xing"
	// FieldLucky holds the string denoting the lucky field in the database.
	FieldLucky = "lucky"
	// FieldExplanation holds the string denoting the explanation field in the database.
	FieldExplanation = "explanation"
	// FieldComment holds the string denoting the comment field in the database.
	FieldComment = "comment"
	// FieldNeedFix holds the string denoting the need_fix field in the database.
	FieldNeedFix = "need_fix"
	// Table holds the table name of the ncharacter in the database.
	Table = "n_character"
)

// Columns holds all SQL columns for ncharacter fields.
var Columns = []string{
	FieldID,
	FieldPinYin,
	FieldChar,
	FieldCharStroke,
	FieldRadical,
	FieldRadicalStroke,
	FieldIsRegular,
	FieldIsSimplified,
	FieldSimplifiedID,
	FieldIsTraditional,
	FieldTraditionalID,
	FieldIsKangXi,
	FieldKangXiID,
	FieldKangXiStroke,
	FieldIsVariant,
	FieldVariantID,
	FieldIsScience,
	FieldScienceStroke,
	FieldWuXing,
	FieldLucky,
	FieldExplanation,
	FieldComment,
	FieldNeedFix,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// DefaultPinYin holds the default value on creation for the "pin_yin" field.
	DefaultPinYin []string
	// DefaultChar holds the default value on creation for the "char" field.
	DefaultChar string
	// DefaultCharStroke holds the default value on creation for the "char_stroke" field.
	DefaultCharStroke int
	// DefaultRadical holds the default value on creation for the "radical" field.
	DefaultRadical string
	// DefaultRadicalStroke holds the default value on creation for the "radical_stroke" field.
	DefaultRadicalStroke int
	// DefaultIsRegular holds the default value on creation for the "is_regular" field.
	DefaultIsRegular bool
	// DefaultIsSimplified holds the default value on creation for the "is_simplified" field.
	DefaultIsSimplified bool
	// DefaultSimplifiedID holds the default value on creation for the "simplified_id" field.
	DefaultSimplifiedID []int
	// DefaultIsTraditional holds the default value on creation for the "is_traditional" field.
	DefaultIsTraditional bool
	// DefaultTraditionalID holds the default value on creation for the "traditional_id" field.
	DefaultTraditionalID []int
	// DefaultIsKangXi holds the default value on creation for the "is_kang_xi" field.
	DefaultIsKangXi bool
	// DefaultKangXiID holds the default value on creation for the "kang_xi_id" field.
	DefaultKangXiID []int
	// DefaultKangXiStroke holds the default value on creation for the "kang_xi_stroke" field.
	DefaultKangXiStroke int
	// DefaultIsVariant holds the default value on creation for the "is_variant" field.
	DefaultIsVariant bool
	// DefaultVariantID holds the default value on creation for the "variant_id" field.
	DefaultVariantID []int
	// DefaultIsScience holds the default value on creation for the "is_science" field.
	DefaultIsScience bool
	// DefaultScienceStroke holds the default value on creation for the "science_stroke" field.
	DefaultScienceStroke int
	// DefaultWuXing holds the default value on creation for the "wu_xing" field.
	DefaultWuXing string
	// DefaultLucky holds the default value on creation for the "lucky" field.
	DefaultLucky string
	// DefaultExplanation holds the default value on creation for the "explanation" field.
	DefaultExplanation string
	// DefaultComment holds the default value on creation for the "comment" field.
	DefaultComment []string
	// DefaultNeedFix holds the default value on creation for the "need_fix" field.
	DefaultNeedFix bool
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID int
)
