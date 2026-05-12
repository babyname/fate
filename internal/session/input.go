package session

import (
	"time"

	"github.com/babyname/fate/ent"
	"github.com/babyname/fate/internal/naming"
)

type Input struct {
	Last   [2]string
	Born   time.Time
	Sex    naming.Sex
	output *Output
}

func (i *Input) Output() *Output {
	if i.output == nil {
		b := naming.ParseNameBasicFromInput(i.Last, i.Born, i.Sex)
		i.output = &Output{
			basic: b,
			cache: NewCache(),
			name:  make(chan naming.FirstName, 128),
		}
	}
	return i.output
}

type Output struct {
	basic *naming.NameBasic
	cache FilterCache
	name  chan naming.FirstName
}

func (o *Output) Basic() *naming.NameBasic {
	return o.basic
}

func (o *Output) SetLastName(ln [2]*ent.Character) {
	o.basic.LastName = ln
}

func (o *Output) ResetNextName() {
	o.cache.SetCount(0)
}

func (o *Output) NextName() (naming.Name, bool) {
	fn, ok := o.cache.Next()
	if ok {
		return naming.Name{
			NameBasic: o.basic,
			FirstName: fn,
		}, true
	}
	return naming.Name{}, false
}

func (o *Output) Filter(s string) int {
	return len(o.cache.Filter(s))
}

func (o *Output) Total() int {
	return o.cache.Len()
}

func (o *Output) SetCacheFilter(filterCache *PutFilter) {
	o.cache.SetFilter(filterCache)
}
