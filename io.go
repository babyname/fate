package fate

import (
	"time"

	"github.com/babyname/fate/ent"
)

type Sex int //girl:0,boy:1

type Input struct {
	Last   [2]string
	Born   time.Time
	Sex    Sex
	output *Output
}

func (i *Input) Output() *Output {
	if i.output == nil {
		b := parseNameBasicFromInput(i)
		i.output = &Output{
			basic: b,
			names: Cache{},
			name:  make(chan FirstName, 128),
		}
	}
	return i.output
}

type Output struct {
	basic *NameBasic
	count int
	names Cache
	name  chan FirstName
}

func (o *Output) Basic() *NameBasic {
	return o.basic
}

func (o *Output) put(name FirstName) {
	o.names.Put(name)
}

func (o *Output) SetLastName(ln [2]*ent.Character) {
	o.basic.LastName = ln
}

func (o *Output) ResetNext() {
	o.count = 0
}

func (o *Output) NextName() (Name, bool) {
	if o.count < o.names.Len() {
		fn, ok := o.names.GetOne(o.count)
		if ok {
			o.count++
			return Name{
				NameBasic: o.basic,
				FirstName: fn,
			}, true
		}
	}
	return Name{}, false
}

func (o *Output) Total() int {
	return o.names.Len()
}

func (o *Output) getLastStroke(filter Filter) [2]int {
	var strokes [2]int
	strokes[0] = filter.GetCharacterStroke(o.Basic().LastName[0])
	if o.Basic().LastName[1] != nil {
		strokes[1] = filter.GetCharacterStroke(o.Basic().LastName[1])
	}
	return strokes
}
