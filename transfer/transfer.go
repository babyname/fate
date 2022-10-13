package transfer

import (
	"context"
	"fmt"

	"github.com/babyname/fate/ent"
)

const LimitMax = 1000

type Transfer interface {
	Start(ctx context.Context) error
}

type transferDatabase struct {
	Source *ent.Client
	Target *ent.Client
	Tables []string
	Limit  int
}

func (t transferDatabase) Start(ctx context.Context) error {
	var err error
	err = t.Target.Schema.Create(ctx)
	if err != nil {
		return err
	}
	for _, table := range t.Tables {
		switch table {
		case "WuGeLucky":
			err = t.transferWuGeLucky(ctx)
		case "Character":
			err = t.transferCharacter(ctx)
		case "WuXing":
			err = t.transferWuXing(ctx)
		}
		if err != nil {
			return err
		}
	}
	fmt.Println("finished")
	if err := t.Source.Close(); err != nil {
		return err
	}
	if err := t.Target.Close(); err != nil {
		return err
	}
	return nil
}

func (t transferDatabase) transferWuGeLucky(ctx context.Context) error {
	c, err := t.Source.WuGeLucky.Query().Count(ctx)
	if err != nil {
		return err
	}
	if c == 0 {
		return nil
	}

	for i := 0; i < c; i += t.Limit {
		luckies, err := t.Source.WuGeLucky.Query().Limit(t.Limit).Offset(i).All(ctx)
		if err != nil {
			return err
		}
		var bluks []*ent.WuGeLuckyCreate
		for x := range luckies {
			lucky := t.Target.WuGeLucky.Create().SetID(luckies[x].ID).SetWuGeLuckyWithOptional(luckies[x])
			bluks = append(bluks, lucky)
			//fmt.Println("insert wugelucky to database:", i, "total", c, "updated", lucky)
		}

		if len(bluks) != 0 {
			saved, err := t.Target.WuGeLucky.CreateBulk(bluks...).Save(ctx)
			if err != nil {
				return err
			}
			fmt.Println("insert wugelucky to database:", i, "total", c, "updated", len(saved))
		}
	}
	return nil
}

func (t transferDatabase) transferCharacter(ctx context.Context) error {
	c, err := t.Source.Character.Query().Count(ctx)
	if err != nil {
		return err
	}
	if c == 0 {
		return nil
	}

	for i := 0; i < c; i += t.Limit {
		characters, err := t.Source.Character.Query().Limit(t.Limit).Offset(i).All(ctx)
		if err != nil {
			return err
		}
		var bluks []*ent.CharacterCreate
		for x := range characters {
			fmt.Println("character is", characters[x].Ch)
			character := t.Target.Character.Create().SetID(characters[x].ID).SetCharacterWithOptional(characters[x])
			bluks = append(bluks, character)
		}

		if len(bluks) != 0 {
			saved, err := t.Target.Character.CreateBulk(bluks...).Save(ctx)
			if err != nil {
				return err
			}
			fmt.Println("insert character to database:", i, "total", c, "updated", len(saved))
		}
	}
	return nil
}

func (t transferDatabase) transferWuXing(ctx context.Context) error {
	c, err := t.Source.WuXing.Query().Count(ctx)
	if err != nil {
		return err
	}
	if c == 0 {
		return nil
	}

	for i := 0; i < c; i += t.Limit {
		wuxings, err := t.Source.WuXing.Query().Limit(t.Limit).Offset(i).All(ctx)
		if err != nil {
			return err
		}
		var bluks []*ent.WuXingCreate
		for x := range wuxings {
			wuxing := t.Target.WuXing.Create().SetID(wuxings[x].ID).SetWuXingWithOptional(wuxings[x])
			bluks = append(bluks, wuxing)
			//fmt.Println("insert wuxing to database:", i, "total", c, "updated", wuxing)
		}

		if len(bluks) != 0 {
			saved, err := t.Target.WuXing.CreateBulk(bluks...).Save(ctx)
			if err != nil {
				return err
			}
			fmt.Println("insert wuxing to database:", i, "total", c, "updated", len(saved))
		}
	}
	return nil
}

func newTransfer(c *DatabaseConfig) (*transferDatabase, error) {
	source, err := c.Source.Database().BuildClient()
	if err != nil {
		return nil, fmt.Errorf("could not open source database: %v", err)
	}
	target, err := c.Target.Database().BuildClient()
	if err != nil {
		return nil, fmt.Errorf("could not open target database: %v", err)
	}
	if c.Limit <= 0 || c.Limit >= LimitMax {
		c.Limit = LimitMax
	}
	return &transferDatabase{
		Tables: c.Tables,
		Source: source,
		Target: target,
		Limit:  c.Limit,
	}, nil
}

func NewTransfer(config *DatabaseConfig) (Transfer, error) {
	return newTransfer(config)
}

var _ Transfer = (*transferDatabase)(nil)
