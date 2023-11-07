// Code generated by ent, DO NOT EDIT.

package ncharacter

import (
	"entgo.io/ent/dialect/sql"
	"github.com/babyname/fate/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id int32) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int32) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int32) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int32) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int32) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int32) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int32) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int32) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int32) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldLTE(FieldID, id))
}

// PinYin applies equality check predicate on the "pin_yin" field. It's identical to PinYinEQ.
func PinYin(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldEQ(FieldPinYin, v))
}

// Ch applies equality check predicate on the "ch" field. It's identical to ChEQ.
func Ch(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldEQ(FieldCh, v))
}

// ChStroke applies equality check predicate on the "ch_stroke" field. It's identical to ChStrokeEQ.
func ChStroke(v int) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldEQ(FieldChStroke, v))
}

// ChType applies equality check predicate on the "ch_type" field. It's identical to ChTypeEQ.
func ChType(v int) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldEQ(FieldChType, v))
}

// Radical applies equality check predicate on the "radical" field. It's identical to RadicalEQ.
func Radical(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldEQ(FieldRadical, v))
}

// RadicalStroke applies equality check predicate on the "radical_stroke" field. It's identical to RadicalStrokeEQ.
func RadicalStroke(v int) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldEQ(FieldRadicalStroke, v))
}

// Relate applies equality check predicate on the "relate" field. It's identical to RelateEQ.
func Relate(v int32) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldEQ(FieldRelate, v))
}

// RelateKangXi applies equality check predicate on the "relate_kang_xi" field. It's identical to RelateKangXiEQ.
func RelateKangXi(v int32) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldEQ(FieldRelateKangXi, v))
}

// RelateTraditional applies equality check predicate on the "relate_traditional" field. It's identical to RelateTraditionalEQ.
func RelateTraditional(v int32) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldEQ(FieldRelateTraditional, v))
}

// IsNameScience applies equality check predicate on the "is_name_science" field. It's identical to IsNameScienceEQ.
func IsNameScience(v bool) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldEQ(FieldIsNameScience, v))
}

// NameScienceChStroke applies equality check predicate on the "name_science_ch_stroke" field. It's identical to NameScienceChStrokeEQ.
func NameScienceChStroke(v int) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldEQ(FieldNameScienceChStroke, v))
}

// IsRegular applies equality check predicate on the "is_regular" field. It's identical to IsRegularEQ.
func IsRegular(v bool) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldEQ(FieldIsRegular, v))
}

// WuXing applies equality check predicate on the "wu_xing" field. It's identical to WuXingEQ.
func WuXing(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldEQ(FieldWuXing, v))
}

// Lucky applies equality check predicate on the "lucky" field. It's identical to LuckyEQ.
func Lucky(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldEQ(FieldLucky, v))
}

// Comment applies equality check predicate on the "comment" field. It's identical to CommentEQ.
func Comment(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldEQ(FieldComment, v))
}

// PinYinEQ applies the EQ predicate on the "pin_yin" field.
func PinYinEQ(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldEQ(FieldPinYin, v))
}

// PinYinNEQ applies the NEQ predicate on the "pin_yin" field.
func PinYinNEQ(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldNEQ(FieldPinYin, v))
}

// PinYinIn applies the In predicate on the "pin_yin" field.
func PinYinIn(vs ...string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldIn(FieldPinYin, vs...))
}

// PinYinNotIn applies the NotIn predicate on the "pin_yin" field.
func PinYinNotIn(vs ...string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldNotIn(FieldPinYin, vs...))
}

// PinYinGT applies the GT predicate on the "pin_yin" field.
func PinYinGT(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldGT(FieldPinYin, v))
}

// PinYinGTE applies the GTE predicate on the "pin_yin" field.
func PinYinGTE(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldGTE(FieldPinYin, v))
}

// PinYinLT applies the LT predicate on the "pin_yin" field.
func PinYinLT(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldLT(FieldPinYin, v))
}

// PinYinLTE applies the LTE predicate on the "pin_yin" field.
func PinYinLTE(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldLTE(FieldPinYin, v))
}

// PinYinContains applies the Contains predicate on the "pin_yin" field.
func PinYinContains(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldContains(FieldPinYin, v))
}

// PinYinHasPrefix applies the HasPrefix predicate on the "pin_yin" field.
func PinYinHasPrefix(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldHasPrefix(FieldPinYin, v))
}

// PinYinHasSuffix applies the HasSuffix predicate on the "pin_yin" field.
func PinYinHasSuffix(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldHasSuffix(FieldPinYin, v))
}

// PinYinEqualFold applies the EqualFold predicate on the "pin_yin" field.
func PinYinEqualFold(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldEqualFold(FieldPinYin, v))
}

// PinYinContainsFold applies the ContainsFold predicate on the "pin_yin" field.
func PinYinContainsFold(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldContainsFold(FieldPinYin, v))
}

// ChEQ applies the EQ predicate on the "ch" field.
func ChEQ(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldEQ(FieldCh, v))
}

// ChNEQ applies the NEQ predicate on the "ch" field.
func ChNEQ(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldNEQ(FieldCh, v))
}

// ChIn applies the In predicate on the "ch" field.
func ChIn(vs ...string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldIn(FieldCh, vs...))
}

// ChNotIn applies the NotIn predicate on the "ch" field.
func ChNotIn(vs ...string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldNotIn(FieldCh, vs...))
}

// ChGT applies the GT predicate on the "ch" field.
func ChGT(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldGT(FieldCh, v))
}

// ChGTE applies the GTE predicate on the "ch" field.
func ChGTE(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldGTE(FieldCh, v))
}

// ChLT applies the LT predicate on the "ch" field.
func ChLT(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldLT(FieldCh, v))
}

// ChLTE applies the LTE predicate on the "ch" field.
func ChLTE(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldLTE(FieldCh, v))
}

// ChContains applies the Contains predicate on the "ch" field.
func ChContains(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldContains(FieldCh, v))
}

// ChHasPrefix applies the HasPrefix predicate on the "ch" field.
func ChHasPrefix(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldHasPrefix(FieldCh, v))
}

// ChHasSuffix applies the HasSuffix predicate on the "ch" field.
func ChHasSuffix(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldHasSuffix(FieldCh, v))
}

// ChEqualFold applies the EqualFold predicate on the "ch" field.
func ChEqualFold(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldEqualFold(FieldCh, v))
}

// ChContainsFold applies the ContainsFold predicate on the "ch" field.
func ChContainsFold(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldContainsFold(FieldCh, v))
}

// ChStrokeEQ applies the EQ predicate on the "ch_stroke" field.
func ChStrokeEQ(v int) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldEQ(FieldChStroke, v))
}

// ChStrokeNEQ applies the NEQ predicate on the "ch_stroke" field.
func ChStrokeNEQ(v int) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldNEQ(FieldChStroke, v))
}

// ChStrokeIn applies the In predicate on the "ch_stroke" field.
func ChStrokeIn(vs ...int) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldIn(FieldChStroke, vs...))
}

// ChStrokeNotIn applies the NotIn predicate on the "ch_stroke" field.
func ChStrokeNotIn(vs ...int) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldNotIn(FieldChStroke, vs...))
}

// ChStrokeGT applies the GT predicate on the "ch_stroke" field.
func ChStrokeGT(v int) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldGT(FieldChStroke, v))
}

// ChStrokeGTE applies the GTE predicate on the "ch_stroke" field.
func ChStrokeGTE(v int) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldGTE(FieldChStroke, v))
}

// ChStrokeLT applies the LT predicate on the "ch_stroke" field.
func ChStrokeLT(v int) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldLT(FieldChStroke, v))
}

// ChStrokeLTE applies the LTE predicate on the "ch_stroke" field.
func ChStrokeLTE(v int) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldLTE(FieldChStroke, v))
}

// ChTypeEQ applies the EQ predicate on the "ch_type" field.
func ChTypeEQ(v int) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldEQ(FieldChType, v))
}

// ChTypeNEQ applies the NEQ predicate on the "ch_type" field.
func ChTypeNEQ(v int) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldNEQ(FieldChType, v))
}

// ChTypeIn applies the In predicate on the "ch_type" field.
func ChTypeIn(vs ...int) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldIn(FieldChType, vs...))
}

// ChTypeNotIn applies the NotIn predicate on the "ch_type" field.
func ChTypeNotIn(vs ...int) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldNotIn(FieldChType, vs...))
}

// ChTypeGT applies the GT predicate on the "ch_type" field.
func ChTypeGT(v int) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldGT(FieldChType, v))
}

// ChTypeGTE applies the GTE predicate on the "ch_type" field.
func ChTypeGTE(v int) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldGTE(FieldChType, v))
}

// ChTypeLT applies the LT predicate on the "ch_type" field.
func ChTypeLT(v int) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldLT(FieldChType, v))
}

// ChTypeLTE applies the LTE predicate on the "ch_type" field.
func ChTypeLTE(v int) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldLTE(FieldChType, v))
}

// RadicalEQ applies the EQ predicate on the "radical" field.
func RadicalEQ(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldEQ(FieldRadical, v))
}

// RadicalNEQ applies the NEQ predicate on the "radical" field.
func RadicalNEQ(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldNEQ(FieldRadical, v))
}

// RadicalIn applies the In predicate on the "radical" field.
func RadicalIn(vs ...string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldIn(FieldRadical, vs...))
}

// RadicalNotIn applies the NotIn predicate on the "radical" field.
func RadicalNotIn(vs ...string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldNotIn(FieldRadical, vs...))
}

// RadicalGT applies the GT predicate on the "radical" field.
func RadicalGT(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldGT(FieldRadical, v))
}

// RadicalGTE applies the GTE predicate on the "radical" field.
func RadicalGTE(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldGTE(FieldRadical, v))
}

// RadicalLT applies the LT predicate on the "radical" field.
func RadicalLT(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldLT(FieldRadical, v))
}

// RadicalLTE applies the LTE predicate on the "radical" field.
func RadicalLTE(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldLTE(FieldRadical, v))
}

// RadicalContains applies the Contains predicate on the "radical" field.
func RadicalContains(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldContains(FieldRadical, v))
}

// RadicalHasPrefix applies the HasPrefix predicate on the "radical" field.
func RadicalHasPrefix(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldHasPrefix(FieldRadical, v))
}

// RadicalHasSuffix applies the HasSuffix predicate on the "radical" field.
func RadicalHasSuffix(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldHasSuffix(FieldRadical, v))
}

// RadicalEqualFold applies the EqualFold predicate on the "radical" field.
func RadicalEqualFold(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldEqualFold(FieldRadical, v))
}

// RadicalContainsFold applies the ContainsFold predicate on the "radical" field.
func RadicalContainsFold(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldContainsFold(FieldRadical, v))
}

// RadicalStrokeEQ applies the EQ predicate on the "radical_stroke" field.
func RadicalStrokeEQ(v int) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldEQ(FieldRadicalStroke, v))
}

// RadicalStrokeNEQ applies the NEQ predicate on the "radical_stroke" field.
func RadicalStrokeNEQ(v int) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldNEQ(FieldRadicalStroke, v))
}

// RadicalStrokeIn applies the In predicate on the "radical_stroke" field.
func RadicalStrokeIn(vs ...int) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldIn(FieldRadicalStroke, vs...))
}

// RadicalStrokeNotIn applies the NotIn predicate on the "radical_stroke" field.
func RadicalStrokeNotIn(vs ...int) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldNotIn(FieldRadicalStroke, vs...))
}

// RadicalStrokeGT applies the GT predicate on the "radical_stroke" field.
func RadicalStrokeGT(v int) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldGT(FieldRadicalStroke, v))
}

// RadicalStrokeGTE applies the GTE predicate on the "radical_stroke" field.
func RadicalStrokeGTE(v int) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldGTE(FieldRadicalStroke, v))
}

// RadicalStrokeLT applies the LT predicate on the "radical_stroke" field.
func RadicalStrokeLT(v int) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldLT(FieldRadicalStroke, v))
}

// RadicalStrokeLTE applies the LTE predicate on the "radical_stroke" field.
func RadicalStrokeLTE(v int) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldLTE(FieldRadicalStroke, v))
}

// RelateEQ applies the EQ predicate on the "relate" field.
func RelateEQ(v int32) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldEQ(FieldRelate, v))
}

// RelateNEQ applies the NEQ predicate on the "relate" field.
func RelateNEQ(v int32) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldNEQ(FieldRelate, v))
}

// RelateIn applies the In predicate on the "relate" field.
func RelateIn(vs ...int32) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldIn(FieldRelate, vs...))
}

// RelateNotIn applies the NotIn predicate on the "relate" field.
func RelateNotIn(vs ...int32) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldNotIn(FieldRelate, vs...))
}

// RelateGT applies the GT predicate on the "relate" field.
func RelateGT(v int32) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldGT(FieldRelate, v))
}

// RelateGTE applies the GTE predicate on the "relate" field.
func RelateGTE(v int32) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldGTE(FieldRelate, v))
}

// RelateLT applies the LT predicate on the "relate" field.
func RelateLT(v int32) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldLT(FieldRelate, v))
}

// RelateLTE applies the LTE predicate on the "relate" field.
func RelateLTE(v int32) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldLTE(FieldRelate, v))
}

// RelateKangXiEQ applies the EQ predicate on the "relate_kang_xi" field.
func RelateKangXiEQ(v int32) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldEQ(FieldRelateKangXi, v))
}

// RelateKangXiNEQ applies the NEQ predicate on the "relate_kang_xi" field.
func RelateKangXiNEQ(v int32) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldNEQ(FieldRelateKangXi, v))
}

// RelateKangXiIn applies the In predicate on the "relate_kang_xi" field.
func RelateKangXiIn(vs ...int32) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldIn(FieldRelateKangXi, vs...))
}

// RelateKangXiNotIn applies the NotIn predicate on the "relate_kang_xi" field.
func RelateKangXiNotIn(vs ...int32) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldNotIn(FieldRelateKangXi, vs...))
}

// RelateKangXiGT applies the GT predicate on the "relate_kang_xi" field.
func RelateKangXiGT(v int32) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldGT(FieldRelateKangXi, v))
}

// RelateKangXiGTE applies the GTE predicate on the "relate_kang_xi" field.
func RelateKangXiGTE(v int32) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldGTE(FieldRelateKangXi, v))
}

// RelateKangXiLT applies the LT predicate on the "relate_kang_xi" field.
func RelateKangXiLT(v int32) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldLT(FieldRelateKangXi, v))
}

// RelateKangXiLTE applies the LTE predicate on the "relate_kang_xi" field.
func RelateKangXiLTE(v int32) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldLTE(FieldRelateKangXi, v))
}

// RelateTraditionalEQ applies the EQ predicate on the "relate_traditional" field.
func RelateTraditionalEQ(v int32) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldEQ(FieldRelateTraditional, v))
}

// RelateTraditionalNEQ applies the NEQ predicate on the "relate_traditional" field.
func RelateTraditionalNEQ(v int32) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldNEQ(FieldRelateTraditional, v))
}

// RelateTraditionalIn applies the In predicate on the "relate_traditional" field.
func RelateTraditionalIn(vs ...int32) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldIn(FieldRelateTraditional, vs...))
}

// RelateTraditionalNotIn applies the NotIn predicate on the "relate_traditional" field.
func RelateTraditionalNotIn(vs ...int32) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldNotIn(FieldRelateTraditional, vs...))
}

// RelateTraditionalGT applies the GT predicate on the "relate_traditional" field.
func RelateTraditionalGT(v int32) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldGT(FieldRelateTraditional, v))
}

// RelateTraditionalGTE applies the GTE predicate on the "relate_traditional" field.
func RelateTraditionalGTE(v int32) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldGTE(FieldRelateTraditional, v))
}

// RelateTraditionalLT applies the LT predicate on the "relate_traditional" field.
func RelateTraditionalLT(v int32) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldLT(FieldRelateTraditional, v))
}

// RelateTraditionalLTE applies the LTE predicate on the "relate_traditional" field.
func RelateTraditionalLTE(v int32) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldLTE(FieldRelateTraditional, v))
}

// IsNameScienceEQ applies the EQ predicate on the "is_name_science" field.
func IsNameScienceEQ(v bool) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldEQ(FieldIsNameScience, v))
}

// IsNameScienceNEQ applies the NEQ predicate on the "is_name_science" field.
func IsNameScienceNEQ(v bool) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldNEQ(FieldIsNameScience, v))
}

// NameScienceChStrokeEQ applies the EQ predicate on the "name_science_ch_stroke" field.
func NameScienceChStrokeEQ(v int) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldEQ(FieldNameScienceChStroke, v))
}

// NameScienceChStrokeNEQ applies the NEQ predicate on the "name_science_ch_stroke" field.
func NameScienceChStrokeNEQ(v int) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldNEQ(FieldNameScienceChStroke, v))
}

// NameScienceChStrokeIn applies the In predicate on the "name_science_ch_stroke" field.
func NameScienceChStrokeIn(vs ...int) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldIn(FieldNameScienceChStroke, vs...))
}

// NameScienceChStrokeNotIn applies the NotIn predicate on the "name_science_ch_stroke" field.
func NameScienceChStrokeNotIn(vs ...int) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldNotIn(FieldNameScienceChStroke, vs...))
}

// NameScienceChStrokeGT applies the GT predicate on the "name_science_ch_stroke" field.
func NameScienceChStrokeGT(v int) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldGT(FieldNameScienceChStroke, v))
}

// NameScienceChStrokeGTE applies the GTE predicate on the "name_science_ch_stroke" field.
func NameScienceChStrokeGTE(v int) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldGTE(FieldNameScienceChStroke, v))
}

// NameScienceChStrokeLT applies the LT predicate on the "name_science_ch_stroke" field.
func NameScienceChStrokeLT(v int) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldLT(FieldNameScienceChStroke, v))
}

// NameScienceChStrokeLTE applies the LTE predicate on the "name_science_ch_stroke" field.
func NameScienceChStrokeLTE(v int) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldLTE(FieldNameScienceChStroke, v))
}

// IsRegularEQ applies the EQ predicate on the "is_regular" field.
func IsRegularEQ(v bool) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldEQ(FieldIsRegular, v))
}

// IsRegularNEQ applies the NEQ predicate on the "is_regular" field.
func IsRegularNEQ(v bool) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldNEQ(FieldIsRegular, v))
}

// WuXingEQ applies the EQ predicate on the "wu_xing" field.
func WuXingEQ(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldEQ(FieldWuXing, v))
}

// WuXingNEQ applies the NEQ predicate on the "wu_xing" field.
func WuXingNEQ(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldNEQ(FieldWuXing, v))
}

// WuXingIn applies the In predicate on the "wu_xing" field.
func WuXingIn(vs ...string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldIn(FieldWuXing, vs...))
}

// WuXingNotIn applies the NotIn predicate on the "wu_xing" field.
func WuXingNotIn(vs ...string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldNotIn(FieldWuXing, vs...))
}

// WuXingGT applies the GT predicate on the "wu_xing" field.
func WuXingGT(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldGT(FieldWuXing, v))
}

// WuXingGTE applies the GTE predicate on the "wu_xing" field.
func WuXingGTE(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldGTE(FieldWuXing, v))
}

// WuXingLT applies the LT predicate on the "wu_xing" field.
func WuXingLT(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldLT(FieldWuXing, v))
}

// WuXingLTE applies the LTE predicate on the "wu_xing" field.
func WuXingLTE(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldLTE(FieldWuXing, v))
}

// WuXingContains applies the Contains predicate on the "wu_xing" field.
func WuXingContains(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldContains(FieldWuXing, v))
}

// WuXingHasPrefix applies the HasPrefix predicate on the "wu_xing" field.
func WuXingHasPrefix(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldHasPrefix(FieldWuXing, v))
}

// WuXingHasSuffix applies the HasSuffix predicate on the "wu_xing" field.
func WuXingHasSuffix(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldHasSuffix(FieldWuXing, v))
}

// WuXingEqualFold applies the EqualFold predicate on the "wu_xing" field.
func WuXingEqualFold(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldEqualFold(FieldWuXing, v))
}

// WuXingContainsFold applies the ContainsFold predicate on the "wu_xing" field.
func WuXingContainsFold(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldContainsFold(FieldWuXing, v))
}

// LuckyEQ applies the EQ predicate on the "lucky" field.
func LuckyEQ(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldEQ(FieldLucky, v))
}

// LuckyNEQ applies the NEQ predicate on the "lucky" field.
func LuckyNEQ(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldNEQ(FieldLucky, v))
}

// LuckyIn applies the In predicate on the "lucky" field.
func LuckyIn(vs ...string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldIn(FieldLucky, vs...))
}

// LuckyNotIn applies the NotIn predicate on the "lucky" field.
func LuckyNotIn(vs ...string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldNotIn(FieldLucky, vs...))
}

// LuckyGT applies the GT predicate on the "lucky" field.
func LuckyGT(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldGT(FieldLucky, v))
}

// LuckyGTE applies the GTE predicate on the "lucky" field.
func LuckyGTE(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldGTE(FieldLucky, v))
}

// LuckyLT applies the LT predicate on the "lucky" field.
func LuckyLT(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldLT(FieldLucky, v))
}

// LuckyLTE applies the LTE predicate on the "lucky" field.
func LuckyLTE(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldLTE(FieldLucky, v))
}

// LuckyContains applies the Contains predicate on the "lucky" field.
func LuckyContains(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldContains(FieldLucky, v))
}

// LuckyHasPrefix applies the HasPrefix predicate on the "lucky" field.
func LuckyHasPrefix(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldHasPrefix(FieldLucky, v))
}

// LuckyHasSuffix applies the HasSuffix predicate on the "lucky" field.
func LuckyHasSuffix(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldHasSuffix(FieldLucky, v))
}

// LuckyEqualFold applies the EqualFold predicate on the "lucky" field.
func LuckyEqualFold(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldEqualFold(FieldLucky, v))
}

// LuckyContainsFold applies the ContainsFold predicate on the "lucky" field.
func LuckyContainsFold(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldContainsFold(FieldLucky, v))
}

// CommentEQ applies the EQ predicate on the "comment" field.
func CommentEQ(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldEQ(FieldComment, v))
}

// CommentNEQ applies the NEQ predicate on the "comment" field.
func CommentNEQ(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldNEQ(FieldComment, v))
}

// CommentIn applies the In predicate on the "comment" field.
func CommentIn(vs ...string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldIn(FieldComment, vs...))
}

// CommentNotIn applies the NotIn predicate on the "comment" field.
func CommentNotIn(vs ...string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldNotIn(FieldComment, vs...))
}

// CommentGT applies the GT predicate on the "comment" field.
func CommentGT(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldGT(FieldComment, v))
}

// CommentGTE applies the GTE predicate on the "comment" field.
func CommentGTE(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldGTE(FieldComment, v))
}

// CommentLT applies the LT predicate on the "comment" field.
func CommentLT(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldLT(FieldComment, v))
}

// CommentLTE applies the LTE predicate on the "comment" field.
func CommentLTE(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldLTE(FieldComment, v))
}

// CommentContains applies the Contains predicate on the "comment" field.
func CommentContains(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldContains(FieldComment, v))
}

// CommentHasPrefix applies the HasPrefix predicate on the "comment" field.
func CommentHasPrefix(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldHasPrefix(FieldComment, v))
}

// CommentHasSuffix applies the HasSuffix predicate on the "comment" field.
func CommentHasSuffix(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldHasSuffix(FieldComment, v))
}

// CommentEqualFold applies the EqualFold predicate on the "comment" field.
func CommentEqualFold(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldEqualFold(FieldComment, v))
}

// CommentContainsFold applies the ContainsFold predicate on the "comment" field.
func CommentContainsFold(v string) predicate.NCharacter {
	return predicate.NCharacter(sql.FieldContainsFold(FieldComment, v))
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.NCharacter) predicate.NCharacter {
	return predicate.NCharacter(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.NCharacter) predicate.NCharacter {
	return predicate.NCharacter(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.NCharacter) predicate.NCharacter {
	return predicate.NCharacter(sql.NotPredicates(p))
}