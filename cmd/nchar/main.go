package main

import (
	"context"
	"fmt"

	_ "github.com/sqlite3ent/sqlite3"

	"github.com/babyname/fate/config"
	"github.com/babyname/fate/database"
	"github.com/babyname/fate/ent"
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
	total := 0

	per := 500
	//var cs []*ent.Characters
	//for i := 0; i < count; i += per {
	//	fmt.Println("update character", "offset", i)
	//	cs, err = client.Characters.Query().Offset(i).Limit(per).All(ctx)
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

	count, err = client.NCharacter.Query().Count(ctx)
	if err != nil {
		panic(err)
	}

	if count == 0 {
		return
	}
	var chs []*ent.NCharacter
	total = 0
	for i := 0; i < count; i += per {
		fmt.Println("update character", "offset", i)
		chs, err = client.NCharacter.Query().Offset(i).Limit(per).All(ctx)
		if err != nil {
			fmt.Println("found error on", "offset", i, "limit", per, "error", err)
			continue
		}
		var up *ent.NCharacterUpdateOne
		for i := range chs {
			total++
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
				fmt.Println("update character went something wrong", "error:", err)
			}
		}
	}
	fmt.Println("update character total", total)
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

	//source.LoadWord("word.json", func(w source.Word) bool {
	//	if len(w.Word) == 0 {
	//		return true
	//	}
	//	exist, err := client.NCharacter.Get(ctx, int([]rune(w.Word)[0]))
	//	if err != nil {
	//		fmt.Println("character not exist", "id", int([]rune(w.Word)[0]), "char", w.Word, "pinyin", w.Pinyin)
	//		s, _ := strconv.ParseInt(w.Strokes, 10, 32)
	//		_, err = client.NCharacter.Create().
	//			SetID(int(int([]rune(w.Word)[0]))).
	//			SetChar(string([]rune{rune(int([]rune(w.Word)[0]))})).
	//			SetCharStroke(int(s)).
	//			SetRadicalStroke(0).
	//			SetRadical(w.Radicals).
	//			SetPinYin([]string{w.Pinyin}).SetNeedFix(true).Save(ctx)
	//		if err != nil {
	//			fmt.Println("something went wrong", "error", err)
	//		}
	//		return true
	//	}
	//	up := exist.Update()
	//	need := false
	//	if exist.Radical != w.Radicals {
	//		need = true
	//		fmt.Println("radical is not a equal", "char", w.Word, "source", exist.Radical, "target", w.Radicals)
	//		if exist.Radical == "" {
	//			up.SetRadical(w.Radicals)
	//		}
	//	}
	//	s, _ := strconv.ParseInt(w.Strokes, 10, 32)
	//	if exist.CharStroke != int(s) {
	//		need = true
	//		fmt.Println("strokes is not a equal", "char", w.Word, "source", exist.CharStroke, "target", w.Strokes)
	//		if exist.CharStroke == 0 {
	//			up.SetCharStroke(int(s))
	//		}
	//	}
	//	_, _ = up.SetNeedFix(need).Save(ctx)
	//	return true
	//})

	//err = source.LoadCharDetailJSON("char_detail.json", func(ch source.Character) bool {
	//	total++
	//	if len(ch.Char) == 0 {
	//		return true
	//	}
	//	var pinyin []string
	//	for _, pron := range ch.Pronunciations {
	//		//fmt.Println("pinyin", pron.Pinyin)
	//		//pron.Pinyin
	//		//for _, exp := range pron.Explanations {
	//
	//		//}
	//		pinyin = append(pinyin, pron.Pinyin)
	//	}
	//	exist, err := client.NCharacter.Get(ctx, int([]rune(ch.Char)[0]))
	//	if err != nil {
	//		fmt.Println("character not exist", "id", int([]rune(ch.Char)[0]), "char", ch.Char)
	//
	//		_, err = client.NCharacter.Create().
	//			SetID(int(int([]rune(ch.Char)[0]))).
	//			SetChar(string([]rune{rune(int([]rune(ch.Char)[0]))})).
	//			//SetCharStroke(int(phone.Strokes)).
	//			SetRadicalStroke(0).
	//			//SetRadical(phone.).
	//			SetPinYin(pinyin).SetNeedFix(true).Save(ctx)
	//		if err != nil {
	//			fmt.Println("something went wrong", "error", err)
	//		}
	//		return true
	//	}
	//	up := exist.Update()
	//	need := false
	//	if len(exist.PinYin) != len(pinyin) {
	//		need = true
	//		fmt.Println("pinyin is not a equal", "char", ch.Char, "source", exist.PinYin, "pinyin", pinyin)
	//	}
	//	_, _ = up.SetNeedFix(need).Save(ctx)
	//	return true
	//})
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("char detail count", total)

	total = 0
	err = source.LoadKangXiChar("kangxi-strokecount.csv", func(kx source.KangXi) bool {
		total++
		if len(kx.Character) == 0 {
			return true
		}
		exist, err := client.NCharacter.Get(ctx, int([]rune(kx.Character)[0]))
		if err != nil {
			fmt.Println("character not exist", "id", int([]rune(kx.Character)[0]), "char", kx.Character, "stroke", kx.Strokes)
			_, err = client.NCharacter.Create().
				SetID(int(int([]rune(kx.Character)[0]))).
				SetChar(string([]rune{rune(int([]rune(kx.Character)[0]))})).
				SetCharStroke(int(kx.Strokes)).
				//SetRadicalStroke(0).
				//SetRadical(w.Radicals).
				//SetPinYin([]string{w.Pinyin})
				SetKangXiStroke(kx.Strokes).
				SetIsKangXi(true).
				SetNeedFix(true).Save(ctx)
			if err != nil {
				fmt.Println("something went wrong", "error", err)
			}
			return true
		}
		up := exist.Update()
		up.SetIsKangXi(true)
		up.SetKangXiStroke(kx.Strokes)
		need := false
		//fmt.Println("kangxi", "cp", kx.CodePoint, "id", kx.Value)
		if exist.CharStroke != kx.Strokes {
			need = true
			fmt.Println("strokes is not a equal", "char", kx.Character, "source", exist.CharStroke, "target", kx.Strokes)
			if exist.CharStroke == 0 {
				up.SetCharStroke(kx.Strokes)
			} else {
				up.SetKangXiStroke(kx.Strokes)
			}
		}
		_, _ = up.SetNeedFix(need).Save(ctx)
		return true
	})
	fmt.Println("kangxi char count", total)

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
