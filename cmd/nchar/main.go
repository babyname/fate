package main

import (
	"context"
	"fmt"

	"github.com/babyname/fate/config"
	"github.com/babyname/fate/database"
	"github.com/babyname/fate/ent"
	"github.com/babyname/fate/ent/schema"
)

func main() {
	cli := database.New(config.DefaultConfig().Database)
	client, err := cli.Client()
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	count, err := client.Character.Query().Count(ctx)
	if err != nil {
		panic(err)
	}

	if count == 0 {
		return
	}

	per := 500
	var cs []*ent.Character
	//var err error
	for i := 0; i < count; i += per {
		cs, err = client.Character.Query().Offset(i).Limit(per).All(ctx)
		if err != nil {
			fmt.Println("found error on", "offset", i, "limit", per, "error", err)
			continue
		}
		for csi := range cs {
			if cs[csi].Ch == "" {
				continue
			}
			var vc []string
			if len(cs[csi].VariantCharacter) > 0 {
				vctmp := []rune(cs[csi].VariantCharacter)
				for i := 0; i < len(vctmp); i++ {
					vc = append(vc, string(vctmp[i]))
				}
			}
			chid := int32([]rune(cs[csi].Ch)[0])
			kxid := int32(0)
			if len([]rune(cs[csi].KangXi)) > 0 {
				kxid = int32([]rune(cs[csi].KangXi)[0])
			}
			tcid := int32(0)
			if len([]rune(cs[csi].TraditionalCharacter)) > 0 {
				tcid = int32([]rune(cs[csi].TraditionalCharacter)[0])
			}
			nc := ent.NCharacter{
				ID:                  chid,
				PinYin:              cs[csi].PinYin,
				Ch:                  cs[csi].Ch,
				ChStroke:            cs[csi].Stroke,
				ChType:              schema.CharTypeSimple,
				Radical:             cs[csi].Radical,
				RadicalStroke:       cs[csi].RadicalStroke,
				Relate:              0,
				RelateKangXi:        kxid,
				RelateTraditional:   tcid,
				RelateVariant:       vc,
				IsNameScience:       cs[csi].NameScience,
				NameScienceChStroke: cs[csi].ScienceStroke,
				IsRegular:           cs[csi].Regular,
				WuXing:              cs[csi].WuXing,
				Lucky:               cs[csi].Lucky,
				Comment:             cs[csi].Comment,
			}
			_, err := client.NCharacter.Create().SetNCharacterWithOptional(&nc).Save(ctx)
			if err != nil {
				continue
			}
			if len([]rune(cs[csi].KangXi)) > 0 {
				kxnc := ent.NCharacter{
					ID:                  []rune(cs[csi].KangXi)[0],
					PinYin:              cs[csi].PinYin,
					Ch:                  cs[csi].KangXi,
					ChStroke:            cs[csi].KangXiStroke,
					ChType:              schema.CharTypeKangXi,
					Radical:             cs[csi].Radical,
					RadicalStroke:       cs[csi].RadicalStroke,
					Relate:              chid,
					RelateKangXi:        kxid,
					RelateTraditional:   tcid,
					RelateVariant:       vc,
					IsNameScience:       cs[csi].NameScience,
					NameScienceChStroke: cs[csi].ScienceStroke,
					IsRegular:           cs[csi].Regular,
					WuXing:              cs[csi].WuXing,
					Lucky:               cs[csi].Lucky,
					Comment:             cs[csi].Comment,
				}
				_, err := client.NCharacter.Create().SetNCharacterWithOptional(&kxnc).Save(ctx)
				if err != nil {
					continue
				}
			}
			if len([]rune(cs[csi].TraditionalCharacter)) > 0 {
				tc := ent.NCharacter{
					ID:                  []rune(cs[csi].TraditionalCharacter)[0],
					PinYin:              cs[csi].PinYin,
					Ch:                  cs[csi].TraditionalCharacter,
					ChStroke:            cs[csi].TraditionalTotalStroke,
					ChType:              schema.CharTypeKangXi,
					Radical:             cs[csi].TraditionalRadical,
					RadicalStroke:       cs[csi].TraditionalRadicalStroke,
					Relate:              chid,
					RelateKangXi:        tcid,
					RelateTraditional:   kxid,
					RelateVariant:       vc,
					IsNameScience:       cs[csi].NameScience,
					NameScienceChStroke: cs[csi].ScienceStroke,
					IsRegular:           cs[csi].Regular,
					WuXing:              cs[csi].WuXing,
					Lucky:               cs[csi].Lucky,
					Comment:             cs[csi].Comment,
				}
				_, err := client.NCharacter.Create().SetNCharacterWithOptional(&tc).Save(ctx)
				if err != nil {
					continue
				}
			}
		}
	}

}
