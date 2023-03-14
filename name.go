package fate

import (
	"fmt"

	"github.com/babyname/fate/ent"
	"github.com/godcong/chronos/v2"
)

type NameBase struct {
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
	*NameBase
	FirstName
}

func (n Name) String() string {
	if n.LastName[1] != nil {
		return fmt.Sprintf("%v%v %v%v", n.LastName[0].Ch, n.LastName[1].Ch, n.FirstName[0].Ch, n.FirstName[1].Ch)
	}
	return fmt.Sprintf("%v %v%v", n.LastName[0].Ch, n.FirstName[0].Ch, n.FirstName[1].Ch)
}
