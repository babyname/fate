package fate

import (
	"fmt"

	"github.com/babyname/fate/ent"
	"github.com/godcong/chronos/v2"
)

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
		return fmt.Sprintf("%v%v %v%v", n.LastName[0].Ch, n.LastName[1].Ch, n.FirstName[0].Ch, n.FirstName[1].Ch)
	}
	return fmt.Sprintf("%v %v%v", n.LastName[0].Ch, n.FirstName[0].Ch, n.FirstName[1].Ch)
}

func (n Name) Strokes() string {
	if n.LastName[1] == nil {
		return fmt.Sprintf("%v,%v,%v", n.LastName[0].ScienceStroke, n.FirstName[0].ScienceStroke, n.FirstName[1].ScienceStroke)
	}
	return fmt.Sprintf("%v,%v,%v,%v", n.LastName[0].ScienceStroke, n.LastName[1].ScienceStroke, n.FirstName[0].ScienceStroke, n.FirstName[1].ScienceStroke)
}

func parseNameBasicFromInput(input *Input) *NameBasic {
	b := chronos.ParseTime(input.Born)
	return &NameBasic{
		LastName:  [2]*ent.Character{},
		Born:      b,
		EightChar: b.Lunar().GetEightChar(),
		Sex:       input.Sex,
	}
}
