package source

import (
	"fmt"

	"golang.org/x/net/context"

	"github.com/babyname/fate/ent"
	"github.com/babyname/fate/ent/ncharacter"
)

func NewCharacterUpdate(client *ent.Client, character *ent.Character) *CharacterUpdate {
	return &CharacterUpdate{
		cli:     client,
		Source:  character,
		Simps:   nil,
		Trads:   nil,
		KangXis: nil,
	}
}

type CharacterUpdate struct {
	cli     *ent.Client
	Source  *ent.Character
	Simps   []int
	Trads   []int
	KangXis []int
	Varies  []int
}

func (c *CharacterUpdate) parseID() {
	c.Simps = append(c.Simps, int([]rune(c.Source.Ch)[0]))
	if len([]rune(c.Source.KangXi)) > 0 {
		c.KangXis = append(c.KangXis, int([]rune(c.Source.KangXi)[0]))
	}
	for _, s := range c.Source.TraditionalCharacter {
		c.Trads = append(c.Trads, int([]rune(s)[0]))
	}
	for _, s := range c.Source.VariantCharacter {
		c.Varies = append(c.Varies, int([]rune(s)[0]))
	}
}

func (c *CharacterUpdate) insertSimple(ctx context.Context) error {
	if len(c.Simps) == 0 {
		return nil
	}
	exist, err := c.cli.NCharacter.Query().Where(ncharacter.ID(c.Simps[0])).Exist(ctx)
	if err != nil {
		return err
	}
	if exist {
		fmt.Println("simple characters already exist", "id", c.Simps[0])
		return nil
	}
	nc := &ent.NCharacter{
		ID:            c.Simps[0],
		PinYin:        c.Source.PinYin,
		Char:          c.Source.Ch,
		CharStroke:    c.Source.Stroke,
		Radical:       c.Source.Radical,
		RadicalStroke: c.Source.RadicalStroke,
		IsRegular:     c.Source.Regular,
		IsSimplified:  true,
		SimplifiedID:  c.simpIDs(),
		IsTraditional: false,
		TraditionalID: c.TradIds(),
		IsKangXi:      false,
		KangXiID:      c.kangxiIDs(),
		IsVariant:     false,
		VariantID:     c.variantIDs(),
		IsScience:     c.Source.NameScience,
		ScienceStroke: c.Source.ScienceStroke,
		WuXing:        c.Source.WuXing,
		Lucky:         c.Source.Lucky,
		Explanation:   "",
		Comment:       c.Source.Comment,
	}
	c.SetSimpleChar(nc)
	nc, err = c.cli.NCharacter.Create().SetNCharacterWithOptional(nc).SetID(c.Simps[0]).Save(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (c *CharacterUpdate) insertKangxi(ctx context.Context) error {
	if len(c.KangXis) == 0 {
		return nil
	}
	exist, err := c.cli.NCharacter.Query().Where(ncharacter.ID(c.KangXis[0])).Exist(ctx)
	if err != nil {
		return err
	}
	if exist {
		fmt.Println("kangxi characters already exist", "id", c.KangXis[0])
		return nil
	}
	nc := &ent.NCharacter{
		ID:            c.KangXis[0],
		PinYin:        c.Source.PinYin,
		CharStroke:    c.Source.KangXiStroke,
		Radical:       c.Source.Radical,
		RadicalStroke: c.Source.RadicalStroke,
		IsRegular:     false,
		IsSimplified:  false,
		SimplifiedID:  c.simpIDs(),
		IsTraditional: false,
		TraditionalID: c.TradIds(),
		IsKangXi:      true,
		KangXiID:      c.kangxiIDs(),
		IsVariant:     false,
		VariantID:     c.variantIDs(),
		IsScience:     c.Source.NameScience,
		ScienceStroke: c.Source.ScienceStroke,
		WuXing:        c.Source.WuXing,
		Lucky:         c.Source.Lucky,
		Explanation:   "",
		Comment:       c.Source.Comment,
	}
	c.SetKangXiChar(nc)
	nc, err = c.cli.NCharacter.Create().SetNCharacterWithOptional(nc).SetID(c.KangXis[0]).Save(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (c *CharacterUpdate) insertTraditional(ctx context.Context) error {
	if len(c.Trads) == 0 {
		return nil
	}
	for _, trad := range c.Trads {
		exist, err := c.cli.NCharacter.Query().Where(ncharacter.ID(trad)).Exist(ctx)
		if err != nil {
			return err
		}
		if exist {
			fmt.Println("trad characters already exist", "id", trad)
			continue
		}
		nc := &ent.NCharacter{
			ID:            trad,
			PinYin:        c.Source.PinYin,
			Char:          string([]rune{rune(trad)}),
			CharStroke:    c.Source.TraditionalTotalStroke,
			Radical:       c.Source.TraditionalRadical,
			RadicalStroke: c.Source.TraditionalRadicalStroke,
			IsRegular:     false,
			IsSimplified:  false,
			SimplifiedID:  c.simpIDs(),
			IsTraditional: true,
			TraditionalID: c.TradIds(),
			IsKangXi:      false,
			KangXiID:      c.kangxiIDs(),
			IsVariant:     false,
			VariantID:     c.variantIDs(),
			IsScience:     c.Source.NameScience,
			ScienceStroke: c.Source.ScienceStroke,
			WuXing:        c.Source.WuXing,
			Lucky:         c.Source.Lucky,
			Explanation:   "",
			Comment:       c.Source.Comment,
		}
		c.SetTradChar(nc)
		nc, err = c.cli.NCharacter.Create().SetNCharacterWithOptional(nc).SetID(trad).Save(ctx)
		if err != nil {
			return err
		}
		//c.Trads = append(c.Trads, nc.ID)
	}
	return nil
}

func (c *CharacterUpdate) insertVariant(ctx context.Context) error {
	if len(c.Varies) == 0 {
		return nil
	}
	for _, vari := range c.Varies {
		exist, err := c.cli.NCharacter.Query().Where(ncharacter.ID(vari)).Exist(ctx)
		if err != nil {
			return err
		}
		if exist {
			fmt.Println("vari characters already exist", "id", vari)
			continue
		}
		nc := &ent.NCharacter{
			ID:            vari,
			PinYin:        c.Source.PinYin,
			Char:          string([]rune{rune(vari)}),
			IsRegular:     false,
			IsSimplified:  false,
			SimplifiedID:  c.simpIDs(),
			IsTraditional: false,
			TraditionalID: c.TradIds(),
			IsKangXi:      false,
			KangXiID:      c.kangxiIDs(),
			IsVariant:     true,
			VariantID:     c.variantIDs(),
			IsScience:     c.Source.NameScience,
			ScienceStroke: c.Source.ScienceStroke,
			WuXing:        c.Source.WuXing,
			Lucky:         c.Source.Lucky,
			Explanation:   "",
			Comment:       c.Source.Comment,
		}
		nc, err = c.cli.NCharacter.Create().SetNCharacterWithOptional(nc).SetID(vari).Save(ctx)
		if err != nil {
			return err
		}
		//c.Trads = append(c.Trads, nc.ID)
	}
	return nil
}

func (c *CharacterUpdate) simpIDs() (ids []int) {
	return c.Simps
}

func (c *CharacterUpdate) TradIds() (ids []int) {
	return c.Trads
}

func (c *CharacterUpdate) kangxiIDs() (ids []int) {
	return c.KangXis
}

func (c *CharacterUpdate) variantIDs() (ids []int) {
	return c.Varies
}

func (c *CharacterUpdate) SetKangXiChar(ch *ent.NCharacter) {
	if c.Source.IsKangXi && c.Source.Ch == "" {
		ch.Char = c.Source.KangXi
	} else {
		ch.Char = c.Source.Ch
	}

	if c.Source.Stroke == 0 && c.Source.Stroke != c.Source.KangXiStroke {
		ch.CharStroke = c.Source.KangXiStroke
	} else {
		ch.CharStroke = c.Source.Stroke
	}
}

func (c *CharacterUpdate) SetTradChar(nc *ent.NCharacter) {
	if c.Source.RadicalStroke == 0 && c.Source.RadicalStroke != c.Source.TraditionalRadicalStroke {
		nc.RadicalStroke = c.Source.TraditionalRadicalStroke
		nc.Radical = c.Source.TraditionalRadical
		nc.CharStroke = c.Source.TraditionalTotalStroke
	} else {
		nc.RadicalStroke = c.Source.RadicalStroke
		nc.Radical = c.Source.Radical
		nc.CharStroke = c.Source.Stroke
	}
}

func (c *CharacterUpdate) SetSimpleChar(nc *ent.NCharacter) {
	if c.Source.RadicalStroke == 0 && c.Source.RadicalStroke != c.Source.SimpleRadicalStroke {
		nc.RadicalStroke = c.Source.SimpleRadicalStroke
		nc.Radical = c.Source.SimpleRadical
		nc.CharStroke = c.Source.SimpleTotalStroke
	} else {
		nc.RadicalStroke = c.Source.RadicalStroke
		nc.Radical = c.Source.Radical
		nc.CharStroke = c.Source.Stroke
	}
}

func (c *CharacterUpdate) Update(ctx context.Context) error {
	c.parseID()
	err := c.insertSimple(ctx)
	if err != nil {
		fmt.Println("failed to insert simple", "error:", err)
		return err
	}
	err = c.insertKangxi(ctx)
	if err != nil {
		fmt.Println("failed to insert kangxi", "error:", err)
		return err
	}
	err = c.insertTraditional(ctx)
	if err != nil {
		fmt.Println("failed to insert translation", "error:", err)
		return err
	}
	err = c.insertVariant(ctx)
	if err != nil {
		fmt.Println("failed to insert variant", "error:", err)
		return err
	}
	return nil
}
