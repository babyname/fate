package main

import (
	"context"
	"fmt"
	"strconv"

	_ "github.com/sqlite3ent/sqlite3"

	"github.com/babyname/fate/config"
	"github.com/babyname/fate/database"
	"github.com/babyname/fate/source"
)

func main() {
	cli := database.New(config.DefaultConfig().Database)
	client, err := cli.Client()
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fmt.Println("update table")
	err = client.Schema.Create(ctx)
	if err != nil {
		panic(err)
	}

	count, err := client.Character.Query().Count(ctx)
	if err != nil {
		panic(err)
	}

	if count == 0 {
		return
	}

	//per := 500
	//var cs []*ent.Character
	//for i := 0; i < count; i += per {
	//	fmt.Println("update character", "offset", i)
	//	cs, err = client.Character.Query().Offset(i).Limit(per).All(ctx)
	//	if err != nil {
	//		fmt.Println("found error on", "offset", i, "limit", per, "error", err)
	//		continue
	//	}
	//	var cu *source.CharacterUpdate
	//	for csi := range cs {
	//		if cs[csi].Ch == "" {
	//			continue
	//		}
	//		cu = source.NewCharacterUpdate(client, cs[csi])
	//		err := cu.Update(ctx)
	//		if err != nil {
	//			continue
	//		}
	//	}
	//}
	//
	//source.LoadPinYin("zdic.txt", func(yin *source.PinYin) bool {
	//	exist, err := client.NCharacter.Get(ctx, int(yin.ID))
	//	if err != nil {
	//		fmt.Println("pinyin not exist", "id", yin.ID, "char", yin.Char, "pinyin", yin.Pinyin)
	//		_, _ = client.NCharacter.Create().SetID(int(yin.ID)).SetChar(string([]rune{rune(yin.ID)})).SetPinYin(yin.Pinyin).SetNeedFix(true).Save(ctx)
	//		return true
	//	}
	//	if len(exist.PinYin) != len(yin.Pinyin) {
	//		fmt.Println("pinyin is not a equal", "char", yin.Char, "source", exist.PinYin, "update", yin.Pinyin)
	//	}
	//	_, err = exist.Update().SetPinYin(mergePinYin(exist.PinYin, yin.Pinyin)).Save(ctx)
	//	if err != nil {
	//		return true
	//	}
	//	return true
	//})

	source.LoadWord("word.json", func(w source.Word) bool {
		if len(w.Word) == 0 {
			return true
		}
		exist, err := client.NCharacter.Get(ctx, int([]rune(w.Word)[0]))
		if err != nil {
			fmt.Println("character not exist", "id", int([]rune(w.Word)[0]), "char", w.Word, "pinyin", w.Pinyin)
			s, _ := strconv.ParseInt(w.Strokes, 10, 32)
			_, err = client.NCharacter.Create().
				SetID(int(int([]rune(w.Word)[0]))).
				SetChar(string([]rune{rune(int([]rune(w.Word)[0]))})).
				SetCharStroke(int(s)).
				SetRadicalStroke(0).
				SetRadical(w.Radicals).
				SetPinYin([]string{w.Pinyin}).SetNeedFix(true).Save(ctx)
			if err != nil {
				fmt.Println("something went wrong", "error", err)
			}
			return true
		}
		up := exist.Update()
		need := false
		if exist.Radical != w.Radicals {
			need = true
			fmt.Println("radical is not a equal", "char", w.Word, "source", exist.Radical, "target", w.Radicals)
			if exist.Radical == "" {
				up.SetRadical(w.Radicals)
			}
		}
		s, _ := strconv.ParseInt(w.Strokes, 10, 32)
		if exist.CharStroke != int(s) {
			need = true
			fmt.Println("strokes is not a equal", "char", w.Word, "source", exist.CharStroke, "target", w.Strokes)
			if exist.CharStroke == 0 {
				up.SetCharStroke(int(s))
			}
		}
		_, _ = up.SetNeedFix(need).Save(ctx)
		return true
	})
}

func mergePinYin(source, target []string) []string {
	tmp := make(map[string]struct{})
	for _, s := range source {
		tmp[s] = struct{}{}
	}
	for _, s := range target {
		tmp[s] = struct{}{}
	}
	var ret []string
	for s := range tmp {
		ret = append(ret, s)
	}
	return ret
}
