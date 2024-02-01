package scripts

import (
	"golang.org/x/net/context"

	"github.com/babyname/fate/ent"
)

func FixEmptyArray(ctx context.Context, client *ent.Client) error {
	count, err := client.NCharacter.Query().Count(ctx)
	if err != nil {
		return err
	}

	if count == 0 {
		return nil
	}
	var chs []*ent.NCharacter
	for i := 0; i < count; i += perLimit {
		chs, err = client.NCharacter.Query().Offset(i).Limit(perLimit).All(ctx)
		if err != nil {
			return err
		}
		var up *ent.NCharacterUpdateOne
		for i := range chs {
			up = chs[i].Update()
			if len(chs[i].PinYin) == 0 {
				up.SetPinYin([]string{})
			}
			if len(chs[i].SimplifiedID) == 0 {
				up.SetSimplifiedID([]int{})
			}
			if len(chs[i].TraditionalID) == 0 {
				up.SetTraditionalID([]int{})
			}
			if len(chs[i].KangXiID) == 0 {
				up.SetKangXiID([]int{})
			}
			if len(chs[i].VariantID) == 0 {
				up.SetVariantID([]int{})
			}
			if len(chs[i].Comment) == 0 {
				up.SetComment([]string{})
			}
			_, err = up.Save(ctx)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
