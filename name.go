package fate

import (
	"fmt"

	"github.com/godcong/chronos/v2"

	"github.com/babyname/fate/ent"
)

type FirstName [2]*ent.Character
type LastName [2]*ent.Character

type BasicInfo struct {
	LastName  [2]*ent.Character
	Born      chronos.Calendar
	EightChar chronos.EightChar
	Sex       Sex
}

type Names struct {
	*BasicInfo
	firsts []FirstName
}

type Name struct {
	*BasicInfo
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

func parseNameBasicFromInput(input *Input) *BasicInfo {
	b := chronos.ParseTime(input.Born)
	return &BasicInfo{
		LastName:  [2]*ent.Character{},
		Born:      b,
		EightChar: b.Lunar().GetEightChar(),
		Sex:       input.Sex,
	}
}
