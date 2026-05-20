package naming

import (
	"fmt"
	"time"

	"github.com/babyname/chronos/v2"
	"github.com/babyname/fate/ent"
)

// Sex 表示性别类型。
type Sex int

const (
	// SexBoy 表示男性。
	SexBoy Sex = 1
	// SexGirl 表示女性。
	SexGirl Sex = 0
)

// NameBasic 表示姓名基础信息，包含姓氏、出生时间、八字和性别。
type NameBasic struct {
	LastName  [2]*ent.Character
	Born      chronos.Calendar
	EightChar chronos.EightChar
	Sex       Sex
}

// FirstName 表示名字部分，由两个汉字组成。
type FirstName [2]*ent.Character

// LastName 表示姓氏部分，由一到两个汉字组成。
type LastName [2]*ent.Character

// NameDetail 表示名字的详细信息。
type NameDetail struct {
	First [2]*ent.Character
}

// Name 表示完整姓名，包含基础信息和名字。
type Name struct {
	*NameBasic
	FirstName
}

func (n Name) String() string {
	if n.LastName[1] != nil {
		return fmt.Sprintf("%v%v %v%v", n.LastName[0].Char, n.LastName[1].Char, n.FirstName[0].Char, n.FirstName[1].Char)
	}
	return fmt.Sprintf("%v %v%v", n.LastName[0].Char, n.FirstName[0].Char, n.FirstName[1].Char)
}

// Strokes 返回姓名各字的科学笔画数，以逗号分隔。
func (n Name) Strokes() string {
	if n.LastName[1] == nil {
		return fmt.Sprintf("%v,%v,%v", n.LastName[0].ScienceStroke, n.FirstName[0].ScienceStroke, n.FirstName[1].ScienceStroke)
	}
	return fmt.Sprintf("%v,%v,%v,%v", n.LastName[0].ScienceStroke, n.LastName[1].ScienceStroke, n.FirstName[0].ScienceStroke, n.FirstName[1].ScienceStroke)
}

// ParseNameBasicFromInput 根据输入参数构造姓名基础信息。
func ParseNameBasicFromInput(_ [2]string, born time.Time, sex Sex) *NameBasic {
	b := chronos.ParseTime(born)
	return &NameBasic{
		LastName:  [2]*ent.Character{},
		Born:      b,
		EightChar: b.Lunar().GetEightChar(),
		Sex:       sex,
	}
}
