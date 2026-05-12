package naming

import (
	"fmt"
	"time"

	"github.com/babyname/fate/ent"
	"github.com/godcong/chronos/v2"
)

type Sex int

type NameBasic struct {
	LastName  [2]*ent.Character
	Born      chronos.Calendar
	EightChar chronos.EightChar
	Sex       Sex
}

type FirstName [2]*ent.Character

type LastName [2]*ent.Character

type NameDetail struct {
	First [2]*ent.Character
}

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

func (n Name) Strokes() string {
	if n.LastName[1] == nil {
		return fmt.Sprintf("%v,%v,%v", n.LastName[0].ScienceStroke, n.FirstName[0].ScienceStroke, n.FirstName[1].ScienceStroke)
	}
	return fmt.Sprintf("%v,%v,%v,%v", n.LastName[0].ScienceStroke, n.LastName[1].ScienceStroke, n.FirstName[0].ScienceStroke, n.FirstName[1].ScienceStroke)
}

func ParseNameBasicFromInput(last [2]string, born time.Time, sex Sex) *NameBasic {
	b := chronos.ParseTime(born)
	return &NameBasic{
		LastName:  [2]*ent.Character{},
		Born:      b,
		EightChar: b.Lunar().GetEightChar(),
		Sex:       sex,
	}
}
