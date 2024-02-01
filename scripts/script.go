package scripts

import (
	"strconv"

	"golang.org/x/net/context"

	"github.com/babyname/fate/ent"
	"github.com/babyname/fate/log"
)

type Script struct {
	c *ent.Client
}

func (s *Script) OutputCharNeedFix(ctx context.Context) error {
	return FindCharNeedFix(ctx, s.c, func(fix NeedFix) bool {
		log.Logger("script", "fix", fix)
		return true
	})
}

func (s *Script) TransferOldToNChar(ctx context.Context) error {
	count, err := s.c.Character.Query().Count(ctx)
	if err != nil {
		return err
	}

	if count == 0 {
		return nil
	}
	var cs []*ent.Character
	for i := 0; i < count; i += perLimit {
		log.Logger("script").Info("update character", "offset", i)
		cs, err = s.c.Character.Query().Offset(i).Limit(perLimit).All(ctx)
		if err != nil {
			return err
		}
		var cu *CharacterUpdate
		for csi := range cs {
			if cs[csi].Ch == "" {
				continue
			}
			cu = NewCharacterUpdate(s.c, cs[csi])
			err := cu.Update(ctx)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Script) FixEmpty(ctx context.Context) error {
	return FixEmptyArray(ctx, s.c)
}

func (s *Script) OutputCharDetailJSON(ctx context.Context, file string) error {
	return LoadCharDetailJSON(file, func(ch Character) bool {
		if len(ch.Char) == 0 {
			return true
		}
		var pinyin []string
		for _, pron := range ch.Pronunciations {
			//log.Logger("script").Info("pinyin", pron.Pinyin)
			//pron.Pinyin
			//for _, exp := range pron.Explanations {

			//}
			pinyin = append(pinyin, pron.Pinyin)
		}
		exist, err := s.c.NCharacter.Get(ctx, int([]rune(ch.Char)[0]))
		if err != nil {
			log.Logger("script").Info("character not exist", "id", int([]rune(ch.Char)[0]), "char", ch.Char)

			_, err = s.c.NCharacter.Create().
				SetID(int(int([]rune(ch.Char)[0]))).
				SetChar(string([]rune{rune(int([]rune(ch.Char)[0]))})).
				//SetCharStroke(int(phone.Strokes)).
				SetRadicalStroke(0).
				//SetRadical(phone.).
				SetPinYin(pinyin).SetNeedFix(true).Save(ctx)
			if err != nil {
				log.Logger("script").Info("something went wrong", "error", err)
			}
			return true
		}
		up := exist.Update()
		need := false
		if len(exist.PinYin) != len(pinyin) {
			need = true
			log.Logger("script").Info("pinyin is not a equal", "char", ch.Char, "source", exist.PinYin, "pinyin", pinyin)
		}
		_, _ = up.SetNeedFix(need).Save(ctx)
		return true
	})
}

func (s *Script) OutputWord(ctx context.Context, file string) error {
	return LoadWord("word.json", func(w Word) bool {
		if len(w.Word) == 0 {
			return true
		}
		exist, err := s.c.NCharacter.Get(ctx, int([]rune(w.Word)[0]))
		if err != nil {
			log.Logger("script").Info("character not exist", "id", int([]rune(w.Word)[0]), "char", w.Word, "pinyin", w.Pinyin)
			ss, _ := strconv.ParseInt(w.Strokes, 10, 32)
			_, err = s.c.NCharacter.Create().
				SetID(int(int([]rune(w.Word)[0]))).
				SetChar(string([]rune{rune(int([]rune(w.Word)[0]))})).
				SetCharStroke(int(ss)).
				SetRadicalStroke(0).
				SetRadical(w.Radicals).
				SetPinYin([]string{w.Pinyin}).SetNeedFix(true).Save(ctx)
			if err != nil {
				log.Logger("script").Info("something went wrong", "error", err)
			}
			return true
		}
		up := exist.Update()
		need := false
		if exist.Radical != w.Radicals {
			need = true
			log.Logger("script").Info("radical is not a equal", "char", w.Word, "source", exist.Radical, "target", w.Radicals)
			if exist.Radical == "" {
				up.SetRadical(w.Radicals)
			}
		}
		s, _ := strconv.ParseInt(w.Strokes, 10, 32)
		if exist.CharStroke != int(s) {
			need = true
			log.Logger("script").Info("strokes is not a equal", "char", w.Word, "source", exist.CharStroke, "target", w.Strokes)
			if exist.CharStroke == 0 {
				up.SetCharStroke(int(s))
			}
		}
		_, _ = up.SetNeedFix(need).Save(ctx)
		return true
	})
}

func (s *Script) OutputPinYin(ctx context.Context, file string) error {
	return LoadPinYin(file, func(yin *PinYin) bool {
		exist, err := s.c.NCharacter.Get(ctx, int(yin.ID))
		if err != nil {
			log.Logger("scripts").Info("pinyin not exist", "id", yin.ID, "char", yin.Char, "pinyin", yin.Pinyin)
			_, _ = s.c.NCharacter.Create().SetID(int(yin.ID)).SetChar(string([]rune{rune(yin.ID)})).SetPinYin(yin.Pinyin).SetNeedFix(true).Save(ctx)
			return true
		}
		if len(exist.PinYin) != len(yin.Pinyin) {
			log.Logger("scripts").Info("pinyin is not a equal", "char", yin.Char, "source", exist.PinYin, "update", yin.Pinyin)
		}
		_, err = exist.Update().SetPinYin(mergePinYin(exist.PinYin, yin.Pinyin)).Save(ctx)
		if err != nil {
			return true
		}
		return true
	})
}

func (s *Script) OutputKangXiChar(ctx context.Context, file string) error {
	return LoadKangXiChar(file, func(kx KangXi) bool {
		if len(kx.Character) == 0 {
			return true
		}
		exist, err := s.c.NCharacter.Get(ctx, int([]rune(kx.Character)[0]))
		if err != nil {
			log.Logger("script").Info("character not exist", "id", int([]rune(kx.Character)[0]), "char", kx.Character, "stroke", kx.Strokes)
			_, err = s.c.NCharacter.Create().
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
				log.Logger("script").Error("error saving", "error", err)
			}
			return true
		}
		up := exist.Update()
		up.SetIsKangXi(true)
		up.SetKangXiStroke(kx.Strokes)
		need := false

		if exist.CharStroke != kx.Strokes {
			need = true
			log.Logger("script").Info("strokes is not a equal", "char", kx.Character, "source", exist.CharStroke, "target", kx.Strokes)
			if exist.CharStroke == 0 {
				up.SetCharStroke(kx.Strokes)
			} else {
				up.SetKangXiStroke(kx.Strokes)
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
