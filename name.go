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
	Last  [2]*ent.Character
}

type Name struct {
	*NameBase
	*NameDetail
}

func (n Name) String() string {
	if n.Last[1] != nil {
		return fmt.Sprintf("%v%v %v%v", n.Last[0], n.Last[1], n.First[0], n.First[1])
	}
	return fmt.Sprintf("%v %v%v", n.Last[0], n.First[0], n.First[1])
}
