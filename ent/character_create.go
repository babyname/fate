// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/babyname/fate/ent/character"
)

// CharacterCreate is the builder for creating a Character entity.
type CharacterCreate struct {
	config
	mutation *CharacterMutation
	hooks    []Hook
}

// SetPinYin sets the "pin_yin" field.
func (cc *CharacterCreate) SetPinYin(s string) *CharacterCreate {
	cc.mutation.SetPinYin(s)
	return cc
}

// SetCh sets the "ch" field.
func (cc *CharacterCreate) SetCh(s string) *CharacterCreate {
	cc.mutation.SetCh(s)
	return cc
}

// SetRadical sets the "radical" field.
func (cc *CharacterCreate) SetRadical(s string) *CharacterCreate {
	cc.mutation.SetRadical(s)
	return cc
}

// SetRadicalStroke sets the "radical_stroke" field.
func (cc *CharacterCreate) SetRadicalStroke(i int) *CharacterCreate {
	cc.mutation.SetRadicalStroke(i)
	return cc
}

// SetStroke sets the "stroke" field.
func (cc *CharacterCreate) SetStroke(i int) *CharacterCreate {
	cc.mutation.SetStroke(i)
	return cc
}

// SetIsKangXi sets the "is_kang_xi" field.
func (cc *CharacterCreate) SetIsKangXi(b bool) *CharacterCreate {
	cc.mutation.SetIsKangXi(b)
	return cc
}

// SetKangXi sets the "kang_xi" field.
func (cc *CharacterCreate) SetKangXi(s string) *CharacterCreate {
	cc.mutation.SetKangXi(s)
	return cc
}

// SetKangXiStroke sets the "kang_xi_stroke" field.
func (cc *CharacterCreate) SetKangXiStroke(i int) *CharacterCreate {
	cc.mutation.SetKangXiStroke(i)
	return cc
}

// SetSimpleRadical sets the "simple_radical" field.
func (cc *CharacterCreate) SetSimpleRadical(s string) *CharacterCreate {
	cc.mutation.SetSimpleRadical(s)
	return cc
}

// SetSimpleRadicalStroke sets the "simple_radical_stroke" field.
func (cc *CharacterCreate) SetSimpleRadicalStroke(i int) *CharacterCreate {
	cc.mutation.SetSimpleRadicalStroke(i)
	return cc
}

// SetSimpleTotalStroke sets the "simple_total_stroke" field.
func (cc *CharacterCreate) SetSimpleTotalStroke(i int) *CharacterCreate {
	cc.mutation.SetSimpleTotalStroke(i)
	return cc
}

// SetTraditionalRadical sets the "traditional_radical" field.
func (cc *CharacterCreate) SetTraditionalRadical(s string) *CharacterCreate {
	cc.mutation.SetTraditionalRadical(s)
	return cc
}

// SetTraditionalRadicalStroke sets the "traditional_radical_stroke" field.
func (cc *CharacterCreate) SetTraditionalRadicalStroke(i int) *CharacterCreate {
	cc.mutation.SetTraditionalRadicalStroke(i)
	return cc
}

// SetTraditionalTotalStroke sets the "traditional_total_stroke" field.
func (cc *CharacterCreate) SetTraditionalTotalStroke(i int) *CharacterCreate {
	cc.mutation.SetTraditionalTotalStroke(i)
	return cc
}

// SetNameScience sets the "name_science" field.
func (cc *CharacterCreate) SetNameScience(b bool) *CharacterCreate {
	cc.mutation.SetNameScience(b)
	return cc
}

// SetWuXing sets the "wu_xing" field.
func (cc *CharacterCreate) SetWuXing(s string) *CharacterCreate {
	cc.mutation.SetWuXing(s)
	return cc
}

// SetLucky sets the "lucky" field.
func (cc *CharacterCreate) SetLucky(s string) *CharacterCreate {
	cc.mutation.SetLucky(s)
	return cc
}

// SetRegular sets the "regular" field.
func (cc *CharacterCreate) SetRegular(b bool) *CharacterCreate {
	cc.mutation.SetRegular(b)
	return cc
}

// SetTraditionalCharacter sets the "traditional_character" field.
func (cc *CharacterCreate) SetTraditionalCharacter(s string) *CharacterCreate {
	cc.mutation.SetTraditionalCharacter(s)
	return cc
}

// SetVariantCharacter sets the "variant_character" field.
func (cc *CharacterCreate) SetVariantCharacter(s string) *CharacterCreate {
	cc.mutation.SetVariantCharacter(s)
	return cc
}

// SetComment sets the "comment" field.
func (cc *CharacterCreate) SetComment(s string) *CharacterCreate {
	cc.mutation.SetComment(s)
	return cc
}

// SetScienceStroke sets the "science_stroke" field.
func (cc *CharacterCreate) SetScienceStroke(i int) *CharacterCreate {
	cc.mutation.SetScienceStroke(i)
	return cc
}

// SetID sets the "id" field.
func (cc *CharacterCreate) SetID(s string) *CharacterCreate {
	cc.mutation.SetID(s)
	return cc
}

// Mutation returns the CharacterMutation object of the builder.
func (cc *CharacterCreate) Mutation() *CharacterMutation {
	return cc.mutation
}

// Save creates the Character in the database.
func (cc *CharacterCreate) Save(ctx context.Context) (*Character, error) {
	var (
		err  error
		node *Character
	)
	if len(cc.hooks) == 0 {
		if err = cc.check(); err != nil {
			return nil, err
		}
		node, err = cc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*CharacterMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = cc.check(); err != nil {
				return nil, err
			}
			cc.mutation = mutation
			if node, err = cc.sqlSave(ctx); err != nil {
				return nil, err
			}
			mutation.id = &node.ID
			mutation.done = true
			return node, err
		})
		for i := len(cc.hooks) - 1; i >= 0; i-- {
			if cc.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = cc.hooks[i](mut)
		}
		v, err := mut.Mutate(ctx, cc.mutation)
		if err != nil {
			return nil, err
		}
		nv, ok := v.(*Character)
		if !ok {
			return nil, fmt.Errorf("unexpected node type %T returned from CharacterMutation", v)
		}
		node = nv
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (cc *CharacterCreate) SaveX(ctx context.Context) *Character {
	v, err := cc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (cc *CharacterCreate) Exec(ctx context.Context) error {
	_, err := cc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cc *CharacterCreate) ExecX(ctx context.Context) {
	if err := cc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (cc *CharacterCreate) check() error {
	if _, ok := cc.mutation.PinYin(); !ok {
		return &ValidationError{Name: "pin_yin", err: errors.New(`ent: missing required field "Character.pin_yin"`)}
	}
	if _, ok := cc.mutation.Ch(); !ok {
		return &ValidationError{Name: "ch", err: errors.New(`ent: missing required field "Character.ch"`)}
	}
	if _, ok := cc.mutation.Radical(); !ok {
		return &ValidationError{Name: "radical", err: errors.New(`ent: missing required field "Character.radical"`)}
	}
	if _, ok := cc.mutation.RadicalStroke(); !ok {
		return &ValidationError{Name: "radical_stroke", err: errors.New(`ent: missing required field "Character.radical_stroke"`)}
	}
	if _, ok := cc.mutation.Stroke(); !ok {
		return &ValidationError{Name: "stroke", err: errors.New(`ent: missing required field "Character.stroke"`)}
	}
	if _, ok := cc.mutation.IsKangXi(); !ok {
		return &ValidationError{Name: "is_kang_xi", err: errors.New(`ent: missing required field "Character.is_kang_xi"`)}
	}
	if _, ok := cc.mutation.KangXi(); !ok {
		return &ValidationError{Name: "kang_xi", err: errors.New(`ent: missing required field "Character.kang_xi"`)}
	}
	if _, ok := cc.mutation.KangXiStroke(); !ok {
		return &ValidationError{Name: "kang_xi_stroke", err: errors.New(`ent: missing required field "Character.kang_xi_stroke"`)}
	}
	if _, ok := cc.mutation.SimpleRadical(); !ok {
		return &ValidationError{Name: "simple_radical", err: errors.New(`ent: missing required field "Character.simple_radical"`)}
	}
	if _, ok := cc.mutation.SimpleRadicalStroke(); !ok {
		return &ValidationError{Name: "simple_radical_stroke", err: errors.New(`ent: missing required field "Character.simple_radical_stroke"`)}
	}
	if _, ok := cc.mutation.SimpleTotalStroke(); !ok {
		return &ValidationError{Name: "simple_total_stroke", err: errors.New(`ent: missing required field "Character.simple_total_stroke"`)}
	}
	if _, ok := cc.mutation.TraditionalRadical(); !ok {
		return &ValidationError{Name: "traditional_radical", err: errors.New(`ent: missing required field "Character.traditional_radical"`)}
	}
	if _, ok := cc.mutation.TraditionalRadicalStroke(); !ok {
		return &ValidationError{Name: "traditional_radical_stroke", err: errors.New(`ent: missing required field "Character.traditional_radical_stroke"`)}
	}
	if _, ok := cc.mutation.TraditionalTotalStroke(); !ok {
		return &ValidationError{Name: "traditional_total_stroke", err: errors.New(`ent: missing required field "Character.traditional_total_stroke"`)}
	}
	if _, ok := cc.mutation.NameScience(); !ok {
		return &ValidationError{Name: "name_science", err: errors.New(`ent: missing required field "Character.name_science"`)}
	}
	if _, ok := cc.mutation.WuXing(); !ok {
		return &ValidationError{Name: "wu_xing", err: errors.New(`ent: missing required field "Character.wu_xing"`)}
	}
	if _, ok := cc.mutation.Lucky(); !ok {
		return &ValidationError{Name: "lucky", err: errors.New(`ent: missing required field "Character.lucky"`)}
	}
	if _, ok := cc.mutation.Regular(); !ok {
		return &ValidationError{Name: "regular", err: errors.New(`ent: missing required field "Character.regular"`)}
	}
	if _, ok := cc.mutation.TraditionalCharacter(); !ok {
		return &ValidationError{Name: "traditional_character", err: errors.New(`ent: missing required field "Character.traditional_character"`)}
	}
	if _, ok := cc.mutation.VariantCharacter(); !ok {
		return &ValidationError{Name: "variant_character", err: errors.New(`ent: missing required field "Character.variant_character"`)}
	}
	if _, ok := cc.mutation.Comment(); !ok {
		return &ValidationError{Name: "comment", err: errors.New(`ent: missing required field "Character.comment"`)}
	}
	if _, ok := cc.mutation.ScienceStroke(); !ok {
		return &ValidationError{Name: "science_stroke", err: errors.New(`ent: missing required field "Character.science_stroke"`)}
	}
	return nil
}

func (cc *CharacterCreate) sqlSave(ctx context.Context) (*Character, error) {
	_node, _spec := cc.createSpec()
	if err := sqlgraph.CreateNode(ctx, cc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	if _spec.ID.Value != nil {
		if id, ok := _spec.ID.Value.(string); ok {
			_node.ID = id
		} else {
			return nil, fmt.Errorf("unexpected Character.ID type: %T", _spec.ID.Value)
		}
	}
	return _node, nil
}

func (cc *CharacterCreate) createSpec() (*Character, *sqlgraph.CreateSpec) {
	var (
		_node = &Character{config: cc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: character.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: character.FieldID,
			},
		}
	)
	if id, ok := cc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = id
	}
	if value, ok := cc.mutation.PinYin(); ok {
		_spec.SetField(character.FieldPinYin, field.TypeString, value)
		_node.PinYin = value
	}
	if value, ok := cc.mutation.Ch(); ok {
		_spec.SetField(character.FieldCh, field.TypeString, value)
		_node.Ch = value
	}
	if value, ok := cc.mutation.Radical(); ok {
		_spec.SetField(character.FieldRadical, field.TypeString, value)
		_node.Radical = value
	}
	if value, ok := cc.mutation.RadicalStroke(); ok {
		_spec.SetField(character.FieldRadicalStroke, field.TypeInt, value)
		_node.RadicalStroke = value
	}
	if value, ok := cc.mutation.Stroke(); ok {
		_spec.SetField(character.FieldStroke, field.TypeInt, value)
		_node.Stroke = value
	}
	if value, ok := cc.mutation.IsKangXi(); ok {
		_spec.SetField(character.FieldIsKangXi, field.TypeBool, value)
		_node.IsKangXi = value
	}
	if value, ok := cc.mutation.KangXi(); ok {
		_spec.SetField(character.FieldKangXi, field.TypeString, value)
		_node.KangXi = value
	}
	if value, ok := cc.mutation.KangXiStroke(); ok {
		_spec.SetField(character.FieldKangXiStroke, field.TypeInt, value)
		_node.KangXiStroke = value
	}
	if value, ok := cc.mutation.SimpleRadical(); ok {
		_spec.SetField(character.FieldSimpleRadical, field.TypeString, value)
		_node.SimpleRadical = value
	}
	if value, ok := cc.mutation.SimpleRadicalStroke(); ok {
		_spec.SetField(character.FieldSimpleRadicalStroke, field.TypeInt, value)
		_node.SimpleRadicalStroke = value
	}
	if value, ok := cc.mutation.SimpleTotalStroke(); ok {
		_spec.SetField(character.FieldSimpleTotalStroke, field.TypeInt, value)
		_node.SimpleTotalStroke = value
	}
	if value, ok := cc.mutation.TraditionalRadical(); ok {
		_spec.SetField(character.FieldTraditionalRadical, field.TypeString, value)
		_node.TraditionalRadical = value
	}
	if value, ok := cc.mutation.TraditionalRadicalStroke(); ok {
		_spec.SetField(character.FieldTraditionalRadicalStroke, field.TypeInt, value)
		_node.TraditionalRadicalStroke = value
	}
	if value, ok := cc.mutation.TraditionalTotalStroke(); ok {
		_spec.SetField(character.FieldTraditionalTotalStroke, field.TypeInt, value)
		_node.TraditionalTotalStroke = value
	}
	if value, ok := cc.mutation.NameScience(); ok {
		_spec.SetField(character.FieldNameScience, field.TypeBool, value)
		_node.NameScience = value
	}
	if value, ok := cc.mutation.WuXing(); ok {
		_spec.SetField(character.FieldWuXing, field.TypeString, value)
		_node.WuXing = value
	}
	if value, ok := cc.mutation.Lucky(); ok {
		_spec.SetField(character.FieldLucky, field.TypeString, value)
		_node.Lucky = value
	}
	if value, ok := cc.mutation.Regular(); ok {
		_spec.SetField(character.FieldRegular, field.TypeBool, value)
		_node.Regular = value
	}
	if value, ok := cc.mutation.TraditionalCharacter(); ok {
		_spec.SetField(character.FieldTraditionalCharacter, field.TypeString, value)
		_node.TraditionalCharacter = value
	}
	if value, ok := cc.mutation.VariantCharacter(); ok {
		_spec.SetField(character.FieldVariantCharacter, field.TypeString, value)
		_node.VariantCharacter = value
	}
	if value, ok := cc.mutation.Comment(); ok {
		_spec.SetField(character.FieldComment, field.TypeString, value)
		_node.Comment = value
	}
	if value, ok := cc.mutation.ScienceStroke(); ok {
		_spec.SetField(character.FieldScienceStroke, field.TypeInt, value)
		_node.ScienceStroke = value
	}
	return _node, _spec
}

// CharacterCreateBulk is the builder for creating many Character entities in bulk.
type CharacterCreateBulk struct {
	config
	builders []*CharacterCreate
}

// Save creates the Character entities in the database.
func (ccb *CharacterCreateBulk) Save(ctx context.Context) ([]*Character, error) {
	specs := make([]*sqlgraph.CreateSpec, len(ccb.builders))
	nodes := make([]*Character, len(ccb.builders))
	mutators := make([]Mutator, len(ccb.builders))
	for i := range ccb.builders {
		func(i int, root context.Context) {
			builder := ccb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*CharacterMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				nodes[i], specs[i] = builder.createSpec()
				var err error
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, ccb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, ccb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, ccb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (ccb *CharacterCreateBulk) SaveX(ctx context.Context) []*Character {
	v, err := ccb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ccb *CharacterCreateBulk) Exec(ctx context.Context) error {
	_, err := ccb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ccb *CharacterCreateBulk) ExecX(ctx context.Context) {
	if err := ccb.Exec(ctx); err != nil {
		panic(err)
	}
}