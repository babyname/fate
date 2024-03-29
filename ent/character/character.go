// Code generated by ent, DO NOT EDIT.

package character

const (
	// Label holds the string label denoting the character type in the database.
	Label = "character"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "hash"
	// FieldPinYin holds the string denoting the pin_yin field in the database.
	FieldPinYin = "pin_yin"
	// FieldCh holds the string denoting the ch field in the database.
	FieldCh = "ch"
	// FieldRadical holds the string denoting the radical field in the database.
	FieldRadical = "radical"
	// FieldRadicalStroke holds the string denoting the radical_stroke field in the database.
	FieldRadicalStroke = "radical_stroke"
	// FieldStroke holds the string denoting the stroke field in the database.
	FieldStroke = "stroke"
	// FieldIsKangXi holds the string denoting the is_kang_xi field in the database.
	FieldIsKangXi = "is_kang_xi"
	// FieldKangXi holds the string denoting the kang_xi field in the database.
	FieldKangXi = "kang_xi"
	// FieldKangXiStroke holds the string denoting the kang_xi_stroke field in the database.
	FieldKangXiStroke = "kang_xi_stroke"
	// FieldSimpleRadical holds the string denoting the simple_radical field in the database.
	FieldSimpleRadical = "simple_radical"
	// FieldSimpleRadicalStroke holds the string denoting the simple_radical_stroke field in the database.
	FieldSimpleRadicalStroke = "simple_radical_stroke"
	// FieldSimpleTotalStroke holds the string denoting the simple_total_stroke field in the database.
	FieldSimpleTotalStroke = "simple_total_stroke"
	// FieldTraditionalRadical holds the string denoting the traditional_radical field in the database.
	FieldTraditionalRadical = "traditional_radical"
	// FieldTraditionalRadicalStroke holds the string denoting the traditional_radical_stroke field in the database.
	FieldTraditionalRadicalStroke = "traditional_radical_stroke"
	// FieldTraditionalTotalStroke holds the string denoting the traditional_total_stroke field in the database.
	FieldTraditionalTotalStroke = "traditional_total_stroke"
	// FieldNameScience holds the string denoting the name_science field in the database.
	FieldNameScience = "name_science"
	// FieldWuXing holds the string denoting the wu_xing field in the database.
	FieldWuXing = "wu_xing"
	// FieldLucky holds the string denoting the lucky field in the database.
	FieldLucky = "lucky"
	// FieldRegular holds the string denoting the regular field in the database.
	FieldRegular = "regular"
	// FieldTraditionalCharacter holds the string denoting the traditional_character field in the database.
	FieldTraditionalCharacter = "traditional_character"
	// FieldVariantCharacter holds the string denoting the variant_character field in the database.
	FieldVariantCharacter = "variant_character"
	// FieldComment holds the string denoting the comment field in the database.
	FieldComment = "comment"
	// FieldScienceStroke holds the string denoting the science_stroke field in the database.
	FieldScienceStroke = "science_stroke"
	// Table holds the table name of the character in the database.
	Table = "character"
)

// Columns holds all SQL columns for character fields.
var Columns = []string{
	FieldID,
	FieldPinYin,
	FieldCh,
	FieldRadical,
	FieldRadicalStroke,
	FieldStroke,
	FieldIsKangXi,
	FieldKangXi,
	FieldKangXiStroke,
	FieldSimpleRadical,
	FieldSimpleRadicalStroke,
	FieldSimpleTotalStroke,
	FieldTraditionalRadical,
	FieldTraditionalRadicalStroke,
	FieldTraditionalTotalStroke,
	FieldNameScience,
	FieldWuXing,
	FieldLucky,
	FieldRegular,
	FieldTraditionalCharacter,
	FieldVariantCharacter,
	FieldComment,
	FieldScienceStroke,
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
